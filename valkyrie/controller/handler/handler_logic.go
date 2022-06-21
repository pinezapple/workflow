package handler

/*
func UserUploadFileToMinio(c echo.Context) (erro error) {
	return controller.ExecHandlerForUploads(c, nil, userUploadFileToMinio)
}
*/

// TODO: maybe use later
/*
func UserUploadFileToMinio(c echo.Context) error {
	//logger init
	logger.Info("user upload to server")
	mainConf := core.GetMainConfig()
	ctx := c.Request().Context()
	userid := ctx.Value("UserID").(string)
	// handler header

	// Get minio client
	s3Client := core.GetMinIOClient()
	db := core.GetDBObj()
	bucketDAO := dao.GetBucketDAO()
	fileDAO := dao.GetFileDAO()

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		logger.Errorf("Extract http multipart form upload error: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, &utilsModel.Response{
			Data: nil,
			Error: utilsModel.ResponseError{
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			},
		})

	}
	files := form.File["files"]
	sample := c.QueryParam("sample")
	workflow_uuid := c.QueryParam("workflow")

	var fileList []*model.GetUserUploadFileList
	now := time.Now().Unix()

	for _, file := range files {
		// Source
		src, err := file.Open()
		if err != nil {
			logger.Errorf("Open upload file error: %s", err.Error())
			return c.JSON(http.StatusInternalServerError, &utilsModel.Response{
				Data: nil,
				Error: utilsModel.ResponseError{
					Message: err.Error(),
					Code:    http.StatusInternalServerError,
				},
			})

		}
		defer src.Close()

		if mainConf.HardDiskOnly {
			// Destination
			/*
				_, err := os.Stat(mainConf.InputDirPrefix + "/" + userid + "-" + file.Filename + "-" + time.Now().String())
				if !os.IsNotExist(err) {
					if rewrite == "1" {
						err = utils.DeleteFile(mainConf.InputDirPrefix + "/" + userid + "-" + file.Filename)
						if err != nil {
							return c.JSON(http.StatusInternalServerError, &utilsModel.Response{
								Data: nil,
								Error: utilsModel.ResponseError{
									Message: fmt.Errorf("failed to rewrite").Error(),
									Code:    http.StatusInternalServerError,
								},
							})

						}
					} else {
						return c.JSON(http.StatusInternalServerError, &utilsModel.Response{
							Data: nil,
							Error: utilsModel.ResponseError{
								Message: fmt.Errorf("file is existed").Error(),
								Code:    http.StatusInternalServerError,
							},
						})
					}
				}

			filePath := mainConf.InputDirPrefix + "/" + userid + "/" + strconv.FormatInt(now, 10) + "/" + file.Filename
			// fileUniqueName := userid + "-" + file.Filename + "-" + strconv.FormatInt(time.Now().Unix(), 10)

			if _, err := os.Stat(filepath.Dir(filePath)); os.IsNotExist(err) {
				_ = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
			}
			dst, err := os.Create(filePath)
			if err != nil {
				logger.Errorf("Create upload local file error: %s", err.Error())
				return c.JSON(http.StatusInternalServerError, &utilsModel.Response{
					Data: nil,
					Error: utilsModel.ResponseError{
						Message: err.Error(),
						Code:    http.StatusInternalServerError,
					},
				})

			}
			defer dst.Close()

			// Copy
			if _, err = io.Copy(dst, src); err != nil {
				logger.Errorf("Write upload file error: %s", err.Error())
				return c.JSON(http.StatusInternalServerError, &utilsModel.Response{
					Data: nil,
					Error: utilsModel.ResponseError{
						Message: err.Error(),
						Code:    http.StatusInternalServerError,
					},
				})

			}

			f := &model.File{
				UserID:         userid,
				SampleName:     sample,
				UserUploadName: file.Filename,
				RunID:          -1,
				TaskID:         "",
				Filename:       file.Filename,
				Filesize:       file.Size,
				LocalPath:      filePath,
				Deleted:        false,
				Safe:           true,
				UploadSuccess:  false,
				WorkflowUUID:   workflow_uuid,
				CreatedAt:      time.Now(),
				ExpiredAt:      time.Now().Add(time.Duration(mainConf.ImportantFileTTL) * time.Hour),
			}
			err = fileDAO.SaveFile(ctx, db, f)
			if err != nil {
				logger.Errorf("Save upload file to db error: %s", err.Error())
				return c.JSON(http.StatusInternalServerError, &utilsModel.Response{
					Data: nil,
					Error: utilsModel.ResponseError{
						Message: err.Error(),
						Code:    http.StatusInternalServerError,
					},
				})

			}

			var uploadedFile = &model.GetUserUploadFileList{
				Path:     filePath,
				Filename: f.UserUploadName,
				Filesize: f.Filesize,
			}
			fileList = append(fileList, uploadedFile)

			continue
		}

		bucket, iter, newBucket := core.GetMinioBucket(file.Size)
		if newBucket {
			err := core.MountBucketToDir(bucket, mainConf.FUSEMountpoint+"/"+bucket, mainConf.MinioAuthenFile, mainConf.MinioEndpoint)
			if err != nil {
				core.DeleteFromBucket(iter, file.Size)
				logger.Errorf("Delete file from bucket error: %s", err.Error())
				return c.JSON(http.StatusInternalServerError, &utilsModel.Response{
					Data: nil,
					Error: utilsModel.ResponseError{
						Message: err.Error(),
						Code:    http.StatusInternalServerError,
					},
				})

			}
		}

		// TODO: init file name
		key := aws.String(userid + "/" + file.Filename)

		_, err = s3Client.PutObject(&s3.PutObjectInput{
			Body:   src,
			Bucket: aws.String(bucket),
			Key:    key,
		})
		if err != nil {
			core.DeleteFromBucket(iter, file.Size)
			logger.Errorf("Delete file from bucket error: %s", err.Error())
			return c.JSON(http.StatusInternalServerError, &utilsModel.Response{
				Data: nil,
				Error: utilsModel.ResponseError{
					Message: err.Error(),
					Code:    http.StatusInternalServerError,
				},
			})

		}

		logger.Info(fmt.Sprintf("Successfully created bucket %s and uploaded data with key %s\n", bucket, *key))

		// TODO : save to db
		if newBucket {
			bucketDAO.CreateNewBucket(ctx, db, bucket, file.Size, iter)
		} else {
			bucketDAO.AddToBucket(ctx, db, bucket, file.Size)
		}

		f := &model.File{
			UserID:        userid,
			RunID:         -1,
			TaskID:        "",
			Bucket:        bucket,
			Filename:      file.Filename,
			Filesize:      file.Size,
			Deleted:       false,
			Safe:          true,
			UploadSuccess: true,
			WorkflowUUID:  workflow_uuid,
			CreatedAt:     time.Now(),
			ExpiredAt:     time.Now().Add(time.Duration(mainConf.ImportantFileTTL) * time.Hour),
		}
		err = fileDAO.SaveFile(ctx, db, f)
		if err != nil {
			logger.Errorf("Save upload file to db error: %s", err.Error())
			return c.JSON(http.StatusInternalServerError, &utilsModel.Response{
				Data: nil,
				Error: utilsModel.ResponseError{
					Message: err.Error(),
					Code:    http.StatusInternalServerError,
				},
			})

		}

		// might be used later
		/*
			// Destination
			dst, err := os.Create(file.Filename)
			if err != nil {
				return err
			}
			defer dst.Close()

			// Copy
			if _, err = io.Copy(dst, src); err != nil {
				return err
			}

	}

	var response = model.GetUserUploadSampleResp{
		SampleName: sample,
		UploadTime: time.Now(),
		FileList:   fileList,
	}
	return c.JSON(http.StatusOK, &utilsModel.Response{
		Data: response,
		Error: utilsModel.ResponseError{
			Message: "",
			Code:    http.StatusOK,
		},
	})

}
*/
