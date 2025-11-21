package cwssql

import (
	"context"
	"log"
	"testing"
	"time"

	_ "ariga.io/atlas-provider-gorm/gormschema"
	"gorm.io/gorm"
)

type Account struct {
	WorkId   string `gorm:"type:text;primaryKey"`
	Account  string `gorm:"type:text;primaryKey"`
	Password string `gorm:"type:text;not null"`
	Enable   bool   `gorm:"default:true"`
	BaseTimeModel
}

type AccountRepository Repository[Account]

func (r *AccountRepository) GetAccount(account string) (*Account, error) {
	return r.Get(Eq("Account", account))
}

func (r *AccountRepository) GetAccounts(account string) ([]*Account, error) {
	return r.GetAll(Eq("Account", account))
}

func NewAccountRepository(context context.Context, session *gorm.DB) *AccountRepository {
	return (*AccountRepository)(NewRepository[Account](context, session))
}

func TestSqlRepository(t *testing.T) {
	db, err := NewSQLiteDB("data.db")
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&Account{})
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	repo := NewAccountRepository(ctx, db)
	err = repo.GetGorm().Transaction(func(tx *gorm.DB) error {

		// err = repo.Begin()
		// if err != nil {
		// 	panic(err)
		// }

		account := &Account{
			Password: "XXXX",
			Account:  "test@xxx.com",
			WorkId:   "XXXXXX",
		}

		account2 := &Account{
			Password: "XXXX",
			Account:  "test2@xxx.com",
			WorkId:   "XXXXXX",
		}

		err = repo.Upsert(account)
		if err != nil {
			panic(err)
		}
		log.Println("Upsert account", account.Account)

		err = repo.Upsert(account2)
		if err != nil {
			panic(err)
		}
		log.Println("Upsert account", account2.Account)

		account.Enable = false
		repo.Refresh(account)

		account.Enable = false
		err = repo.Upsert(account)
		if err != nil {
			panic(err)
		}
		log.Println("Upsert account", account.Enable)

		repo.Refresh(account)

		account, err = repo.GetAccount("test@xxx.com")
		if err != nil {
			panic(err)
		}
		log.Printf("Retrieved account: %+v", account)

		account, err = repo.GetAccount("")
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Println("Account not found")
			} else {
				panic(err)
			}
		}
		log.Printf("Retrieved account: %+v", account)

		accounts, err := repo.GetAccounts("")
		if err != nil {
			log.Printf("Error getting all accounts: %v", err)
			return err
		}
		log.Printf("All account: %+v", accounts)

		accounts, err = repo.GetAll(And(Eq("Enable", false).Eq("Account", "test2@xxx.com"), Eq("Account", "test@xxx.com").Eq("Enable", false)))
		if err != nil {
			log.Printf("Error getting all accounts: %v", err)
			return err
		}
		log.Printf("All account: %+v", accounts)

		accounts, err = repo.GetAll(Or(Eq("Enable", false).Eq("Account", "test2@xxx.com"), Eq("Account", "test@xxx.com").Eq("Enable", false)))
		if err != nil {
			log.Printf("Error getting all accounts: %v", err)
			return err
		}
		log.Printf("All account: %+v", accounts)

		accounts, err = repo.GetAll(Between("CreatedAt", time.Now().Add(-time.Hour), time.Now()))
		if err != nil {
			log.Printf("Error getting all accounts: %v", err)
			return err
		}
		log.Printf("All account: %+v", accounts)

		accounts, err = repo.GetAll(Ne("Enable", true))
		if err != nil {
			log.Printf("Error getting all accounts: %v", err)
			return err
		}
		log.Printf("All account: %+v", accounts)

		// err = repo.Commit()
		if err != nil {
			log.Printf("Error committing transaction: %v", err)
			return err
		}

		accounts, err = repo.DeleteAll(Lte("CreatedAt", time.Now()))
		if err != nil {
			log.Printf("Error deleting all accounts: %v", err)
			return err
		}
		log.Printf("All account: %+v", accounts)
		return nil
	})
	// err = repo.Rollback()
	if err != nil {
		log.Printf("Error rolling back transaction: %v", err)
	}
}
