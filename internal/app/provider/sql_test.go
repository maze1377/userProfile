package provider_test

import (
	"context"
	"io"
	"testing"
	"userProfile/internal/app/provider"

	"userProfile/pkg/sql"
	"userProfile/pkg/userProfile"

	"github.com/stretchr/testify/suite"
	"golang.org/x/xerrors"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type SQLProviderTestSuite struct {
	suite.Suite

	provider provider.ClientInfoProvider
}

func TestSQLProviderTestSuite(t *testing.T) {
	suite.Run(t, new(SQLProviderTestSuite))
}

func (s *SQLProviderTestSuite) TestGetClientInfoShouldReturnNotFoundInitially() {
	_, err := s.provider.GetClientInfo(context.Background(), &userProfile.ClientInfoRequest{
		ClientID: "myToken",
	})
	s.True(xerrors.Is(err, provider.ErrNotFound))
}

func (s *SQLProviderTestSuite) TestShouldReturnClientInfoAfterAdd() {
	err := s.provider.RegisterClientInfo(context.Background(), &userProfile.RegisterRequest{
		UserProfile: &userProfile.UserProfile{
			ClientID: "abcd",
		},
	})
	s.Nil(err)
	if err != nil {
		return
	}

	clientInfo, err := s.provider.GetClientInfo(context.Background(), &userProfile.ClientInfoRequest{
		ClientID: "abcd",
	})
	s.Nil(err)
	if err != nil {
		return
	}

	s.Equal("abcd", clientInfo.ClientID)
}

func (s *SQLProviderTestSuite) SetupTest() {
	db, err := sql.GetDatabase(sql.SqliteConfig{
		InMemory: true,
	})
	if err != nil {
		s.FailNow(err.Error(), "unable to instantiate SQLite instance")
		return
	}

	s.provider = provider.NewSQL(db)

	err = s.provider.(sql.Migrate).Migrate()

	if err != nil {
		s.FailNow(err.Error(), "unable to migrate SQLite database")
	}
}

func (s *SQLProviderTestSuite) TearDownTest() {
	err := s.provider.(io.Closer).Close()
	if err != nil {
		s.FailNow(err.Error(), "unable to close SQLite database")
	}
}
