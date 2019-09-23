package db

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	opentracing "github.com/opentracing/opentracing-go"
	//	"github.com/opentracing/opentracing-go/ext"
	"micro_framework/cmd"
)

// MyDB 封装gorm数据库操作
type MyDB struct {
	*gorm.DB
}

//InitDB 初始化数据库连接
func InitDB(args string) (db *MyDB, err error) {
	myDB, err := gorm.Open("mysql", args)
	if err != nil {
		return nil, err
	}
	return &MyDB{
		DB: myDB,
	}, nil
}

// Create insert the value into database
func (s *MyDB) Create(ctx context.Context, value interface{}) *MyDB {
	span := StartSpanFromContext(ctx, fmt.Sprintf("%+v", value))
	defer FinishSpan(span)
	return &MyDB{
		DB: s.DB.Create(value),
	}
}

// Find find records that match given conditions
func (s *MyDB) Find(ctx context.Context, out interface{}, where ...interface{}) *MyDB {
	span := StartSpanFromContext(ctx, fmt.Sprintf("%+v:%+v", out, where))
	defer FinishSpan(span)
	return &MyDB{
		DB: s.DB.Find(out, where),
	}
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (s *MyDB) Save(ctx context.Context, value interface{}) *MyDB {
	span := StartSpanFromContext(ctx, fmt.Sprintf("%+v", value))
	defer FinishSpan(span)
	return &MyDB{
		DB: s.DB.Save(value),
	}
}

// StartSpanFromContext 开启tracing
func StartSpanFromContext(ctx context.Context, comment string) opentracing.Span {
	tracer := opentracing.GlobalTracer()
	log := cmd.GetLog(ctx)
	log.Infof("ctx:%v", ctx)

	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		log.Infof("span from context is nil ctx:%v", ctx)
		return nil
	}
	serverSpan := tracer.StartSpan(
		comment,
		opentracing.ChildOf(span.Context()),
	)
	return serverSpan
}

//FinishSpan tracing完成
func FinishSpan(span opentracing.Span) {
	if span == nil {
		return
	}
	span.Finish()
}
