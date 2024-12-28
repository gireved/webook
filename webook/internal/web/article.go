package web

import (
	"fmt"
	"geektime-basic-go/webook/internal/domain"
	"geektime-basic-go/webook/internal/service"
	"geektime-basic-go/webook/internal/web/jwt"
	"geektime-basic-go/webook/pkg/ginx"
	"geektime-basic-go/webook/pkg/logger"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"time"
)

type ArticleHandler struct {
	svc     service.ArticleService
	l       logger.LoggerV1
	intrSvc service.InteractiveService
	biz     string
}

func NewArticleHandler(l logger.LoggerV1,
	svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		l:   l,
		svc: svc,
		biz: "article",
	}
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")

	//g.PUT("/", h.Edit)
	g.POST("/edit", ginx.WrapBodyAndClaims(h.Edit))
	g.POST("/publish", ginx.WrapBodyAndClaims(h.Publish))
	g.POST("/withdraw", ginx.WrapBodyAndClaims(h.Withdraw))
	// 创作者的查询接口
	g.POST("/list", h.List)
	g.POST("/detail", h.Detail)

	pub := g.Group("/pub")
	pub.GET("/:id", h.PubDetail)
	// 传入一个参数，true 就是点赞, false 就是不点赞
	pub.POST("/like", ginx.WrapBodyAndClaims(h.Like))
	pub.POST("/collect", ginx.WrapBodyAndClaims(h.Collect))
}

func (h *ArticleHandler) Like(c *gin.Context,
	req ArticleLikeReq, uc jwt.UserClaims) (ginx.Result, error) {
	var err error
	if req.Like {
		// 点赞
		err = h.intrSvc.Like(c, h.biz, req.Id, uc.Uid)
	} else {
		// 取消点赞
		err = h.intrSvc.CancelLike(c, h.biz, req.Id, uc.Uid)
	}
	if err != nil {
		return ginx.Result{
			Code: 5, Msg: "系统错误",
		}, err
	}
	return ginx.Result{
		Msg: "OK",
	}, nil
}

func (h *ArticleHandler) Collect(ctx *gin.Context,
	req ArticleCollectReq, uc jwt.UserClaims) (ginx.Result, error) {
	err := h.intrSvc.Collect(ctx, h.biz, req.Id, req.Cid, uc.Uid)
	if err != nil {
		return ginx.Result{
			Code: 5, Msg: "系统错误",
		}, err
	}
	return ginx.Result{
		Msg: "OK",
	}, nil
}

// Edit 接收 Article 输入，返回一个 ID，文章的 ID
func (h *ArticleHandler) Edit(ctx *gin.Context,
	req ArticleEditReq, uc jwt.UserClaims) (ginx.Result, error) {
	id, err := h.svc.Save(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		return ginx.Result{
			Msg: "系统错误",
		}, err
	}
	return ginx.Result{
		Data: id,
	}, nil
}

func (h *ArticleHandler) Publish(ctx *gin.Context,
	req PublishReq,
	uc jwt.UserClaims) (ginx.Result, error) {
	//val, ok := ctx.Get("user")
	//if !ok {
	//	ctx.JSON(http.StatusOK, Result{
	//		Code: 4,
	//		Msg:  "未登录",
	//	})
	//	return
	//}
	id, err := h.svc.Publish(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Uid,
		},
	})
	if err != nil {
		return ginx.Result{
			Msg:  "系统错误",
			Code: 5,
		}, fmt.Errorf("发表文章失败 aid %d, uid %d %w", uc.Uid, req.Id, err)
	}
	return ginx.Result{
		Data: id,
	}, nil
}

func (h *ArticleHandler) Withdraw(ctx *gin.Context,
	req ArticleWithdrawReq, uc jwt.UserClaims) (ginx.Result, error) {
	err := h.svc.Withdraw(ctx, uc.Uid, req.Id)
	if err != nil {
		return ginx.Result{
			Msg:  "系统错误",
			Code: 5,
		}, err
	}
	return ginx.Result{
		Msg: "OK",
	}, nil
}

func (h *ArticleHandler) List(ctx *gin.Context) {
	var page Page
	if err := ctx.Bind(&page); err != nil {
		return
	}
	// 我要不要检测一下？
	uc := ctx.MustGet("user").(jwt.UserClaims)
	arts, err := h.svc.GetByAuthor(ctx, uc.Uid, page.Offset, page.Limit)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("查找文章列表失败",
			logger.Error(err),
			logger.Int("offset", page.Offset),
			logger.Int("limit", page.Limit),
			logger.Int64("uid", uc.Uid))
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: slice.Map[domain.Article, ArticleVo](arts, func(idx int, src domain.Article) ArticleVo {
			return ArticleVo{
				Id:       src.Id,
				Title:    src.Title,
				Abstract: src.Abstract(),

				//Content:  src.Content,
				AuthorId: src.Author.Id,
				// 列表，你不需要
				Status: src.Status.ToUint8(),
				Ctime:  src.Ctime.Format(time.DateTime),
				Utime:  src.Utime.Format(time.DateTime),
			}
		}),
	})
}

func (h *ArticleHandler) Detail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "参数错误",
		})
		h.l.Error("前端输入的ID不对", logger.Error(err))
		return
	}
	usr, ok := ctx.MustGet("user").(jwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("获取用户会话信息失败")
		return
	}
	art, err := h.svc.GetById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("获取文章信息失败", logger.Error(err))
		return
	}

	if art.Author.Id != usr.Uid {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "输入有误",
		})
		h.l.Error("非法访问文章，创作者 ID 不匹配", logger.Int64("uid", usr.Uid))
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Data: ArticleVo{
			Id:      art.Id,
			Title:   art.Title,
			Status:  art.Status.ToUint8(),
			Content: art.Content,
			Ctime:   art.Ctime.Format(time.DateTime),
			Utime:   art.Utime.Format(time.DateTime),
		},
	})

}

func (h *ArticleHandler) PubDetail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "id 参数错误",
			Code: 4,
		})
		h.l.Warn("查询文章失败，id 格式不对",
			logger.String("id", idstr),
			logger.Error(err))
		return
	}

	var eg errgroup.Group
	var art domain.Article
	uc := ctx.MustGet("users").(jwt.UserClaims)
	eg.Go(func() error {
		// 读文章本体
		art, err = h.svc.GetPubById(ctx, id, uc.Uid)
		return err
	})

	var intr domain.Interactive
	eg.Go(func() error {
		// 可以容忍错误
		intr, err = h.intrSvc.Get(ctx, h.biz, id, uc.Uid)
		return err
	})

	err = eg.Wait()
	if err != nil {
		// 代表查询出错
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	// 增加阅读计数
	go func() {
		er := h.intrSvc.IncrReadCnt(ctx, h.biz, art.Id)
		if er != nil {
			h.l.Error("增加阅读计数失败", logger.Int64("aid", art.Id), logger.Error(err))
		}
	}()

	ctx.JSON(http.StatusOK, Result{
		Data: ArticleVo{
			Id:    art.Id,
			Title: art.Title,

			Content:    art.Content,
			AuthorId:   art.Author.Id,
			AuthorName: art.Author.Name,

			Status:     art.Status.ToUint8(),
			Ctime:      art.Ctime.Format(time.DateTime),
			Utime:      art.Utime.Format(time.DateTime),
			Liked:      intr.Liked,
			Collected:  intr.Collected,
			LikeCnt:    intr.LikeCnt,
			ReadCnt:    intr.ReadCnt,
			CollectCnt: intr.CollectCnt,
		},
	})
}
