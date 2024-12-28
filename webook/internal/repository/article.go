package repository

import (
	"context"
	"geektime-basic-go/webook/internal/domain"
	"geektime-basic-go/webook/internal/repository/cache"
	"geektime-basic-go/webook/internal/repository/dao"
	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, uid int64, id int64, status domain.ArticleStatus) error
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPubById(ctx context.Context, id int64) (domain.Article, error)
}

type CachedArticleRepository struct {
	dao      dao.ArticleDAO
	userRepo UserRepository

	readerDAO dao.ArticleReaderDAO
	authorDAO dao.ArticleAuthorDAO

	db *gorm.DB

	cache cache.ArticleCache
}

func (c *CachedArticleRepository) GetPubById(ctx context.Context, id int64) (domain.Article, error) {
	// 读取线上库数据
	art, err := c.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	// 在这组装user,适合单体应用
	usr, err := c.userRepo.FindById(ctx, art.AuthorId)
	res := domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id:   usr.Id,
			Name: usr.Nickname,
		},
		Ctime:  time.UnixMilli(art.Ctime),
		Utime:  time.UnixMilli(art.Utime),
		Status: domain.ArticleStatus(art.Status),
	}
	return res, nil
}

func (c *CachedArticleRepository) GetById(ctx context.Context, id int64) (domain.Article, error) {
	data, err := c.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return c.toDomain(data), nil

}

func (c *CachedArticleRepository) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	// 首先第一步，判定要不要查询缓存
	// 事实上， limit <= 100 都可以查询缓存
	if offset == 0 && limit == 100 {
		//if offset == 0 && limit <= 100 {
		res, err := c.cache.GetFirstPage(ctx, uid)
		if err == nil {
			return res, err
		} else {
			// 要考虑记录日志
			// 缓存未命中，你是可以忽略的
		}
	}
	arts, err := c.dao.GetByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}
	res := slice.Map[dao.Article, domain.Article](arts, func(idx int, src dao.Article) domain.Article {
		return c.toDomain(src)
	})

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if offset == 0 && limit == 100 {
			// 缓存回写失败，不一定是大问题，但有可能是大问题
			err = c.cache.SetFirstPage(ctx, uid, res)
			if err != nil {
				// 记录日志
				// 我需要监控这里
			}
		}
	}()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		c.preCache(ctx, res)
	}()
	return res, nil
}

func (c *CachedArticleRepository) List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	// 首先第一步，判定要不要查询缓存
	// 事实上， limit <= 100 都可以查询缓存
	if offset == 0 && limit == 100 {
		//if offset == 0 && limit <= 100 {
		res, err := c.cache.GetFirstPage(ctx, uid)
		if err == nil {
			return res, err
		} else {
			// 要考虑记录日志
			// 缓存未命中，你是可以忽略的
		}
	}
	arts, err := c.dao.GetByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}
	res := slice.Map[dao.Article, domain.Article](arts, func(idx int, src dao.Article) domain.Article {
		return c.toDomain(src)
	})

	go func() {
		if offset == 0 && limit == 100 {
			// 缓存回写失败，不一定是大问题，但有可能是大问题
			err = c.cache.SetFirstPage(ctx, uid, res)
			if err != nil {
				// 记录日志
				// 我需要监控这里
			}
			c.preCache(ctx, res)
		}
	}()
	return res, nil
}

func (c *CachedArticleRepository) preCache(ctx context.Context, arts []domain.Article) {
	const size = 1024 * 1024
	if len(arts) > 0 && len(arts[0].Content) < size {
		err := c.cache.Set(ctx, arts[0])
		if err != nil {
			// 记录缓存
		}
	}
}

func (c *CachedArticleRepository) SyncStatus(ctx context.Context, uid int64, id int64, status domain.ArticleStatus) error {
	return c.dao.SyncStatus(ctx, uid, id, status.ToUint8())
}

func (c *CachedArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	id, err := c.dao.Sync(ctx, c.toEntity(art))
	if err == nil {
		er := c.cache.DelFirstPage(ctx, art.Author.Id)
		if er != nil {
			// 也要记录日志
		}
	}
	// 在这里尝试，设置缓存
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// 你可以灵活设置过期时间
		user, er := c.userRepo.FindById(ctx, art.Author.Id)
		if er != nil {
			// 要记录日志
			return
		}
		art.Author = domain.Author{
			Id:   user.Id,
			Name: user.Nickname,
		}
		er = c.cache.SetPub(ctx, art)
		if er != nil {
			// 记录日志
		}
	}()
	return id, err
}

func (c *CachedArticleRepository) SyncV2(ctx context.Context, art domain.Article) (int64, error) {
	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	// 防止后面业务panic
	defer tx.Rollback()

	authorDAO := dao.NewArticleGORMAuthorDAO(tx)
	readerDAO := dao.NewArticleGORMReaderDAO(tx)

	artn := c.toEntity(art)
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = authorDAO.Update(ctx, artn)
	} else {
		id, err = authorDAO.Create(ctx, artn)
	}
	if err != nil {
		return 0, err
	}
	artn.Id = id
	err = readerDAO.UpsertV2(ctx, dao.PublishedArticle(artn))
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return id, nil

}

func (c *CachedArticleRepository) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
	artn := c.toEntity(art)
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = c.authorDAO.Update(ctx, artn)
	} else {
		id, err = c.authorDAO.Create(ctx, artn)
	}
	if err != nil {
		return 0, err
	}
	artn.Id = id
	err = c.readerDAO.Upsert(ctx, artn)
	return id, err
}

func (c *CachedArticleRepository) Update(ctx context.Context, art domain.Article) error {
	defer func() {
		// 清空缓存
		c.cache.DelFirstPage(ctx, art.Author.Id)
	}()
	return c.dao.UpdateById(ctx, c.toEntity(art))
}

func (c *CachedArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	defer func() {
		// 清空缓存
		c.cache.DelFirstPage(ctx, art.Author.Id)
	}()
	return c.dao.Insert(ctx, c.toEntity(art))
}

func NewCachedArticleRepository(dao dao.ArticleDAO, userRepo UserRepository,
	cache cache.ArticleCache) ArticleRepository {
	return &CachedArticleRepository{
		dao:      dao,
		userRepo: userRepo,
		cache:    cache,
	}
}

func NewCachedArticleRepositoryV2(
	readerDAO dao.ArticleReaderDAO,
	authorDAO dao.ArticleAuthorDAO) *CachedArticleRepository {
	return &CachedArticleRepository{
		readerDAO: readerDAO,
		authorDAO: authorDAO,
	}
}

func (c *CachedArticleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		//Status:   uint8(art.Status),
		Status: art.Status.ToUint8(),
	}
}

func (c *CachedArticleRepository) toDomain(art dao.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: art.Id,
		},
		Ctime:  time.UnixMilli(art.Ctime),
		Utime:  time.UnixMilli(art.Utime),
		Status: domain.ArticleStatus(art.Status),
	}
}
