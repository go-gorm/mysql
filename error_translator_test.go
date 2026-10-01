package mysql

import (
	"errors"
	"testing"

	"gorm.io/gorm"

	"github.com/go-sql-driver/mysql"
)

func TestDialector_Translate(t *testing.T) {
	normalErr := errors.New("normal error")

	type fields struct {
		Config *Config
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   error
	}{
		{
			name: "it should translate error to ErrDuplicatedKey when the error number is 1062",
			args: args{err: &mysql.MySQLError{Number: uint16(1062)}},
			want: gorm.ErrDuplicatedKey,
		},
		{
			name: "it should translate error to ErrForeignKeyViolated when the error number is 1451",
			args: args{err: &mysql.MySQLError{Number: uint16(1451)}},
			want: gorm.ErrForeignKeyViolated,
		},
		{
			name: "it should translate error to ErrForeignKeyViolated when the error number is 1452",
			args: args{err: &mysql.MySQLError{Number: uint16(1452)}},
			want: gorm.ErrForeignKeyViolated,
		},
		{
			name: "it should not translate the error when the error number is not registered in translated error codes",
			args: args{err: &mysql.MySQLError{Number: uint16(8888)}},
			want: &mysql.MySQLError{Number: uint16(8888)},
		},
		{
			name: "it should not translate the error when the error is not a mysql error",
			args: args{err: normalErr},
			want: normalErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dialector := Dialector{
				Config: tt.fields.Config,
			}
			if err := dialector.Translate(tt.args.err); !errors.Is(err, tt.want) {
				t.Errorf("Translate() got error = %v, want error %v", err, tt.want)
			}
		})
	}
}

func TestDialector_Translate_PreservesOriginalError(t *testing.T) {
	tests := []struct {
		name        string
		mysqlErr    *mysql.MySQLError
		wantGormErr error
	}{
		{
			name: "it should preserve original MySQLError on ErrDuplicatedKey",
			mysqlErr: &mysql.MySQLError{
				Number:  1062,
				Message: "Duplicate entry 'foo@bar.com' for key 'users.idx_users_email'",
			},
			wantGormErr: gorm.ErrDuplicatedKey,
		},
		{
			name: "it should preserve original MySQLError on ErrForeignKeyViolated when deleting a referenced row",
			mysqlErr: &mysql.MySQLError{
				Number:  1451,
				Message: "Cannot delete or update a parent row: a foreign key constraint fails (`test`.`orders`, CONSTRAINT `fk_orders_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`))",
			},
			wantGormErr: gorm.ErrForeignKeyViolated,
		},
		{
			name: "it should preserve original MySQLError on ErrForeignKeyViolated when inserting an orphan row",
			mysqlErr: &mysql.MySQLError{
				Number:  1452,
				Message: "Cannot add or update a child row: a foreign key constraint fails (`test`.`orders`, CONSTRAINT `fk_orders_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`))",
			},
			wantGormErr: gorm.ErrForeignKeyViolated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dialector := Dialector{}
			err := dialector.Translate(tt.mysqlErr)

			if !errors.Is(err, tt.wantGormErr) {
				t.Errorf("errors.Is() expected %v, got %v", tt.wantGormErr, err)
			}

			var mysqlErr *mysql.MySQLError
			if !errors.As(err, &mysqlErr) {
				t.Fatalf("errors.As() failed: original *mysql.MySQLError is not accessible")
			}

			if mysqlErr.Number != tt.mysqlErr.Number {
				t.Errorf("Number expected %d, got %d", tt.mysqlErr.Number, mysqlErr.Number)
			}
			if mysqlErr.Message != tt.mysqlErr.Message {
				t.Errorf("Message expected %q, got %q", tt.mysqlErr.Message, mysqlErr.Message)
			}
		})
	}
}
