package repository_test

import (
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/eafajri/hr-service.git/module/employee/internal/entity"
	"github.com/eafajri/hr-service.git/module/employee/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Test_AuditLogRepositoryImpl_Create(t *testing.T) {
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

	repo := repository.NewAuditLogRepository(gDb)

	type mocked struct {
		mockDBQueryResult driver.Result
		mockDBQueryErr    error
	}

	testCases := []struct {
		name    string
		log     entity.AuditLog
		payload any
		mocked  mocked
		wantErr error
	}{
		{
			name: "Error Invalid DB",
			log:  entity.AuditLog{},
			mocked: mocked{
				mockDBQueryResult: sqlmock.NewResult(0, 0),
				mockDBQueryErr:    gorm.ErrInvalidDB,
			},
			wantErr: gorm.ErrInvalidDB,
		},
		{
			name:    "Success",
			log:     entity.AuditLog{},
			payload: "{test}",
			mocked: mocked{
				mockDBQueryResult: sqlmock.NewResult(0, 0),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock.ExpectBegin()
			query := "INSERT INTO `audit_logs` (`request_id`,`ip_address`,`table_name`,`action`,`target`,`payload`,`created_by`,`created_at`) VALUES (?,?,?,?,?,CAST(? AS JSON),?,?)"
			mock.ExpectExec(query).
				WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
				WillReturnResult(tc.mocked.mockDBQueryResult).
				WillReturnError(tc.mocked.mockDBQueryErr)

			if tc.wantErr != nil {
				mock.ExpectRollback()
			} else {
				mock.ExpectCommit()
			}

			err := repo.Create(tc.log, nil)
			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
