package service

import (
	"context"
	"errors"
	"geektime-basic-go/webook/internal/domain"
	events "geektime-basic-go/webook/internal/events/article"
	"geektime-basic-go/webook/internal/repository"
	"geektime-basic-go/webook/pkg/logger"
	"time"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, uid int64, id int64) error
	List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPubById(ctx context.Context, id int64, uid int64) (domain.Article, error)
}

type articleService struct {
	repo repository.ArticleRepository

	// V1 写法专用
	readerRepo repository.ArticleReaderRepository
	authorRepo repository.ArticleAuthorRepository
	l          logger.LoggerV1
	producer   events.Producer

	ch chan readInfo
}

type readInfo struct {
	uid int64
	aid int64
}

func (a *articleService) GetPubById(ctx context.Context, id int64, uid int64) (domain.Article, error) {
	art, err := a.repo.GetPubById(ctx, id)
	if err != nil {
		go func() {
			// 生产者也考研通过批量来提高性能

			er := a.producer.ProducerReadEvent(ctx, events.ReadEvent{
				Uid: uid,
				Aid: id,
			})
			if er != nil {
				a.l.Error("发送读者阅读事件失败")
			}
		}()
	}

	// 改批量的做法
	go func() {
		a.ch <- readInfo{
			uid: uid,
			aid: id,
		}
	}()

	return art, err
}

func (a *articleService) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	return a.repo.GetByAuthor(ctx, uid, offset, limit)
}

func (a *articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return a.repo.GetById(ctx, id)
}

func (a *articleService) List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	return a.repo.List(ctx, uid, offset, limit)
}

func (a *articleService) Withdraw(ctx context.Context, uid int64, id int64) error {
	return a.repo.SyncStatus(ctx, uid, id, domain.ArticleStatusPrivate)
}

func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	return a.repo.Sync(ctx, art)
}

func (a *articleService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	// 想到这里要先操作制作库
	// 这里操作线上库
	var (
		id  = art.Id
		err error
	)

	if art.Id > 0 {
		err = a.authorRepo.Update(ctx, art)
	} else {
		id, err = a.authorRepo.Create(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	for i := 0; i < 3; i++ {
		// 我可能线上库已经有数据了
		// 也可能没有
		err = a.readerRepo.Save(ctx, art)
		if err != nil {
			// 多接入一些 tracing 的工具
			a.l.Error("保存到制作库成功但是到线上库失败",
				logger.Int64("aid", art.Id),
				logger.Error(err))
		} else {
			return id, nil
		}
	}
	a.l.Error("保存到制作库成功但是到线上库失败，重试耗尽",
		logger.Int64("aid", art.Id),
		logger.Error(err))
	return id, errors.New("保存到线上库失败，重试次数耗尽")
}

func NewArticleServiceV1(
	readerRepo repository.ArticleReaderRepository,
	authorRepo repository.ArticleAuthorRepository, l logger.LoggerV1) *articleService {
	return &articleService{
		readerRepo: readerRepo,
		authorRepo: authorRepo,
		l:          l,
	}
}

func NewArticleService(repo repository.ArticleRepository, producer events.Producer, l logger.LoggerV1) ArticleService {
	return &articleService{
		repo:     repo,
		producer: producer,
		l:        l,
		//ch:       make(chan readInfo, 10),
	}
}

func NewArticleServiceV2(repo repository.ArticleRepository, producer events.Producer, l logger.LoggerV1) ArticleService {
	ch := make(chan readInfo, 10)
	go func() {
		for {
			uids := make([]int64, 0, 10)
			aids := make([]int64, 0, 10)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			for i := 0; i < 10; i++ {
				select {
				case info, ok := <-ch:
					if !ok {
						cancel()
						return
					}
					uids = append(uids, info.uid)
					aids = append(aids, info.aid)
				case <-ctx.Done():
					break
				}
			}
			cancel()
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			producer.ProducerReadEventV1(ctx, events.ReadEventV1{
				Uids: uids,
				Aids: aids,
			})
			cancel()
		}
	}()
	return &articleService{
		repo:     repo,
		producer: producer,
		l:        l,
		ch:       ch,
	}
}

func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusUnpublished
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	}
	return a.repo.Create(ctx, art)
}
