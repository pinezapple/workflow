package dao

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/vfluxus/valkyrie/model"
)

type bucketDAO struct{}

type IbucketDAO interface {
	CreateNewBucket(ctx context.Context, db *gorm.DB, bucket string, size int64, id int) (err error)
	AddToBucket(ctx context.Context, db *gorm.DB, bucket string, filesize int64) (err error)
	DeleteFromBucket(ctx context.Context, db *gorm.DB, bucket string, filesize int64) (err error)
	SelectAllBucket(ctx context.Context, db *gorm.DB) (result []*model.Bucket, err error)
}

var (
	iBucketDAO IbucketDAO = &bucketDAO{}
)

func SetBucketDAO(v IbucketDAO) {
	iBucketDAO = v
	return
}

func GetBucketDAO() (dao IbucketDAO) {
	return iBucketDAO
}

// -----------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------- IMPLEMENT ----------------------------------------------------------

// CreateNewBucket ...
func (b *bucketDAO) CreateNewBucket(ctx context.Context, db *gorm.DB, bucket string, size int64, id int) (err error) {
	var (
		checkBucket = &model.Bucket{}
	)

	// check if exist
	if err := db.Model(model.Bucket{}).WithContext(ctx).Where("name = ?", bucket).Take(checkBucket).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var (
				newBucket = &model.Bucket{
					ID:           id,
					Name:         bucket,
					CurrentCount: 1,
					CurrentSize:  size,
				}
			)

			if err := db.Model(model.Bucket{}).WithContext(ctx).Create(newBucket).Error; err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}

	// update bucket
	checkBucket.CurrentCount++
	checkBucket.CurrentSize += size

	if err := db.Model(model.Bucket{}).WithContext(ctx).Where("name = ?", bucket).Updates(checkBucket).Error; err != nil {
		return err
	}

	return nil
}

// UpdateBucket ...
func (b *bucketDAO) AddToBucket(ctx context.Context, db *gorm.DB, bucket string, filesize int64) (err error) {
	var (
		checkBucket = &model.Bucket{
			Name: bucket,
		}
	)

	// check if exist
	if err := db.Model(model.Bucket{}).WithContext(ctx).Where("name = ?", bucket).Take(checkBucket).Error; err != nil {
		return err
	}

	// update bucket
	checkBucket.CurrentCount++
	checkBucket.CurrentSize += filesize

	if err := db.Model(model.Bucket{}).WithContext(ctx).Where("name = ?", bucket).Updates(checkBucket).Error; err != nil {
		return err
	}

	return nil
}

func (b *bucketDAO) DeleteFromBucket(ctx context.Context, db *gorm.DB, bucket string, filesize int64) (err error) {
	var (
		checkBucket = &model.Bucket{
			Name: bucket,
		}
	)

	// check if exist
	if err := db.Model(model.Bucket{}).WithContext(ctx).Where("name = ?", bucket).Take(checkBucket).Error; err != nil {
		return err
	}

	// update bucket
	checkBucket.CurrentCount--
	checkBucket.CurrentSize -= filesize

	if err := db.Model(model.Bucket{}).WithContext(ctx).Where("name = ?", bucket).Updates(checkBucket).Error; err != nil {
		return err
	}

	return nil
}

func (b *bucketDAO) SelectAllBucket(ctx context.Context, db *gorm.DB) (result []*model.Bucket, err error) {
	if err := db.Model(&model.Bucket{}).WithContext(ctx).Order("id").Find(&result).Error; err != nil {
		return nil, err
	}
	return result, err
}
