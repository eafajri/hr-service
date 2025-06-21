package repository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/eafajri/hr-service.git/module/employee/internal/entity"
	"github.com/eafajri/hr-service.git/module/employee/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Test_UserRepositoryImpl_GetUserByID(t *testing.T) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer db.Close()
	dialector := mysql.New(mysql.Config{
		Conn: db,
	})
	columns := []string{"version"}
	mock.ExpectQuery("SELECT VERSION()").WithArgs().WillReturnRows(
		mock.NewRows(columns).FromCSVString("1"),
	)

	gDb, _ := gorm.Open(dialector, &gorm.Config{})

	repo := repository.NewUserRepository(gDb)

	type mocked struct {
		mockReturnResult *sqlmock.Rows
		mockDBQueryErr   error
	}

	testCases := []struct {
		name    string
		userID  int64
		mocked  mocked
		wantErr error
		wantRes entity.User
	}{
		{
			name:   "Error Invalid DB",
			userID: 1233,
			mocked: mocked{
				mockReturnResult: sqlmock.NewRows([]string{}),
				mockDBQueryErr:   gorm.ErrInvalidDB,
			},
			wantErr: gorm.ErrInvalidDB,
		},
		{
			name:   "Success",
			userID: 1233,
			mocked: mocked{
				mockReturnResult: sqlmock.NewRows([]string{"id"}).AddRow(1233),
			},
			wantRes: entity.User{ID: 1233},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := "SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT ?"
			mock.ExpectQuery(query).
				WillReturnRows(tc.mocked.mockReturnResult).
				WillReturnError(tc.mocked.mockDBQueryErr)

			res, err := repo.GetUserByID(tc.userID)
			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, res, tc.wantRes)
			}
		})
	}
}

func Test_UserRepositoryImpl_GetUserByUsername(t *testing.T) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer db.Close()
	dialector := mysql.New(mysql.Config{
		Conn: db,
	})
	columns := []string{"version"}
	mock.ExpectQuery("SELECT VERSION()").WithArgs().WillReturnRows(
		mock.NewRows(columns).FromCSVString("1"),
	)

	gDb, _ := gorm.Open(dialector, &gorm.Config{})

	repo := repository.NewUserRepository(gDb)

	type mocked struct {
		mockReturnResult *sqlmock.Rows
		mockDBQueryErr   error
	}

	testCases := []struct {
		name     string
		username string
		mocked   mocked
		wantErr  error
		wantRes  entity.User
	}{
		{
			name:     "Error Invalid DB",
			username: "name",
			mocked: mocked{
				mockReturnResult: sqlmock.NewRows([]string{}),
				mockDBQueryErr:   gorm.ErrInvalidDB,
			},
			wantErr: gorm.ErrInvalidDB,
		},
		{
			name:     "Success",
			username: "name",
			mocked: mocked{
				mockReturnResult: sqlmock.NewRows([]string{"id"}).AddRow(1233),
			},
			wantRes: entity.User{ID: 1233},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := "SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT ?"
			mock.ExpectQuery(query).
				WillReturnRows(tc.mocked.mockReturnResult).
				WillReturnError(tc.mocked.mockDBQueryErr)

			res, err := repo.GetUserByUsername(tc.username)
			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, res, tc.wantRes)
			}
		})
	}
}
