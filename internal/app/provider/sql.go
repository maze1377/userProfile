package provider

import (
	"context"
	"encoding/base64"

	basicError "errors"
	"userProfile/pkg/errors"

	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"

	"userProfile/pkg/userProfile"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type sqlProvider struct {
	write  *gorm.DB
	read   *gorm.DB
	number int
}

type userInfo struct {
	gorm.Model

	ClientId string `gorm:"uniqueIndex"`
	Data     string `gorm:"type:text;"`
}

func NewSQL(db *gorm.DB) ClientInfoProvider {
	return sqlProvider{
		write:  db,
		read:   db,
		number: 1,
	}
}

func NewSQLWithReadAndWrite(write, read *gorm.DB) ClientInfoProvider {
	return sqlProvider{
		write:  write,
		read:   read,
		number: 2,
	}
}

func (p sqlProvider) Close() error {
	sqlDB, err := p.write.DB()
	if err != nil {
		log.Fatalln(err)
		return err
	}
	errWrite := sqlDB.Close()
	var errRead error
	if p.number == 2 {
		sqlDB, err := p.read.DB()
		if err != nil {
			log.Fatalln(err)
			return err
		}
		errRead = sqlDB.Close()
	}
	if errWrite != nil {
		return errWrite
	}
	return errRead
}

func (p sqlProvider) GetClientInfo(ctx context.Context, clientInfo *userProfile.ClientInfoRequest) (*userProfile.UserProfile, error) {
	clientInfoInstance := &userInfo{}
	err := p.read.Where("client_id = ?", clientInfo.GetClientID()).First(clientInfoInstance).Error
	if err != nil {
		if basicError.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WrapWithExtra(ErrNotFound, "userInfo not found", map[string]interface{}{
				"clientId": clientInfo.GetClientID(),
			})
		}
		return nil, errors.WrapWithExtra(err, "could not read userInfo from write", map[string]interface{}{
			"clientId": clientInfo.GetClientID(),
		})
	}

	result, err := p.modelToProto(clientInfoInstance)
	if err != nil {
		return nil, errors.WrapWithExtra(err, "could not convert model to proto", map[string]interface{}{
			"clientId": clientInfo.GetClientID(),
		})
	}

	return result, nil
}

func (p sqlProvider) RegisterClientInfo(ctx context.Context, clientInfo *userProfile.RegisterRequest) error {
	modelInstance, err := p.protoToModel(clientInfo.GetUserProfile())
	if err != nil {
		return errors.WrapWithExtra(err, "could not convert proto to model", map[string]interface{}{
			"userInfo": clientInfo.GetUserProfile(),
		})
	}

	err = p.write.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "client_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"data", "deleted_at", "updated_at"}),
	}).Create(modelInstance).Error
	if err != nil {
		return errors.WrapWithExtra(err, "could not add model to write", map[string]interface{}{
			"userInfo": clientInfo.GetUserProfile(),
		})
	}

	return nil
}

func (p sqlProvider) Migrate() error {
	err := p.write.AutoMigrate(&userInfo{})
	if err != nil {
		return err
	}
	err = p.read.AutoMigrate(&userInfo{})
	return err
}

func (p sqlProvider) protoToModel(userProfileProto *userProfile.UserProfile) (*userInfo, error) {
	binaryData, err := proto.Marshal(userProfileProto)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal proto")
	}

	data := base64.StdEncoding.EncodeToString(binaryData)

	return &userInfo{
		ClientId: userProfileProto.GetClientID(),
		Data:     data,
	}, nil
}

func (p sqlProvider) modelToProto(m *userInfo) (*userProfile.UserProfile, error) {
	data, err := base64.StdEncoding.DecodeString(m.Data)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode base64")
	}

	var result userProfile.UserProfile
	err = proto.Unmarshal(data, &result)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal proto")
	}

	return &result, nil
}
