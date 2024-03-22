package mdbx

import (
	"encoding/hex"
	"errors"
	"github.com/RevittConsulting/chain-dev-utils/internal/types"
	"github.com/RevittConsulting/chain-dev-utils/pkg/utils"
	"github.com/erigontech/mdbx-go/mdbx"
	"log"
)

type MDBX struct {
	env *mdbx.Env
}

func New() *MDBX {
	env, err := mdbx.NewEnv()
	if err != nil {
		log.Fatal(err)
	}

	err = env.SetOption(mdbx.OptMaxDB, 200)
	if err != nil {
		log.Fatal(err)
	}

	return &MDBX{
		env: env,
	}
}

func (m *MDBX) Open(path string) error {
	return m.env.Open(path, mdbx.NoTLS|mdbx.Readonly, 0444)
}

func (m *MDBX) Close() error {
	m.env.Close()
	return nil
}

func (m *MDBX) ListBuckets() ([]string, error) {
	txn, err := m.env.BeginTxn(nil, mdbx.Readonly)
	if err != nil {
		return nil, err
	}
	defer txn.Abort()

	dbis, err := txn.ListDBI()
	if err != nil {
		return nil, err
	}
	return dbis, nil
}

func (m *MDBX) CountKeys(bucketName string) (uint64, error) {
	var count uint64

	txn, err := m.env.BeginTxn(nil, mdbx.Readonly)
	if err != nil {
		return 0, err
	}
	defer txn.Abort()

	actualDbi, err := txn.OpenDBISimple(bucketName, 0)
	if err != nil {
		return 0, err
	}
	dbiStat, err := txn.StatDBI(actualDbi)
	if err != nil {
		return 0, err
	}
	count += dbiStat.Entries

	return count, nil
}

func (m *MDBX) FindByKey(bucketName string, key []byte) ([][]byte, error) {
	txn, err := m.env.BeginTxn(nil, mdbx.Readonly)
	if err != nil {
		return nil, err
	}
	defer txn.Abort()

	dbi, err := txn.OpenDBISimple(bucketName, 0)
	if err != nil {
		return nil, err
	}
	defer m.env.CloseDBI(dbi)

	cursor, err := txn.OpenCursor(dbi)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var values [][]byte
	for k, v, err := cursor.Get(nil, nil, mdbx.First); err == nil; k, v, err = cursor.Get(nil, nil, mdbx.Next) {
		if utils.BytesEqual(k, key) {
			values = append(values, v)
		}
	}
	if err != nil && !errors.Is(err, mdbx.NotFound) {
		return nil, err
	}

	return values, nil
}

func (m *MDBX) FindByValue(bucketName string, value []byte) ([][]byte, error) {
	txn, err := m.env.BeginTxn(nil, mdbx.Readonly)
	if err != nil {
		return nil, err
	}
	defer txn.Abort()

	dbi, err := txn.OpenDBISimple(bucketName, 0)
	if err != nil {
		return nil, err
	}
	defer m.env.CloseDBI(dbi)

	cursor, err := txn.OpenCursor(dbi)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var keys [][]byte
	for k, v, err := cursor.Get(nil, nil, mdbx.First); err == nil; k, v, err = cursor.Get(nil, nil, mdbx.Next) {
		if utils.BytesEqual(v, value) {
			keys = append(keys, k)
		}
	}
	if err != nil && !errors.Is(err, mdbx.NotFound) {
		return nil, err
	}

	return keys, nil
}

func (m *MDBX) Read(bucketName string, take, offset uint64) ([]types.KeyValuePair, error) {
	var data []types.KeyValuePair

	txn, err := m.env.BeginTxn(nil, mdbx.Readonly)
	if err != nil {
		return nil, err
	}
	defer txn.Abort()

	dbi, err := txn.OpenDBISimple(bucketName, 0)
	if err != nil {
		return nil, err
	}
	defer m.env.CloseDBI(dbi)

	cursor, err := txn.OpenCursor(dbi)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	k, v, err := cursor.Get(nil, nil, mdbx.First)
	if err != nil {
		return nil, err
	}

	keyCount, err := m.CountKeys(bucketName)
	if err != nil {
		return nil, err
	}

	if take > keyCount {
		take = keyCount
	}

	count := 0
	for ; err == nil; k, v, err = cursor.Get(nil, nil, mdbx.Next) {
		if uint64(count) >= offset && uint64(count) < (offset+take) {
			data = append(data, types.KeyValuePair{
				Key:   k,
				Value: v,
			})
		}
		count++
		if uint64(count) >= (offset + take) {
			break
		}
	}

	if err != nil && !errors.Is(err, mdbx.NotFound) {
		return nil, err
	}

	return data, nil
}

func (m *MDBX) CountKeysOfLength(bucketName string, length uint64) (uint64, []string, error) {
	var count uint64
	var keys []string

	txn, err := m.env.BeginTxn(nil, mdbx.Readonly)
	if err != nil {
		return 0, nil, err
	}
	defer txn.Abort()

	dbi, err := txn.OpenDBISimple(bucketName, 0)
	if err != nil {
		return 0, nil, err
	}
	defer m.env.CloseDBI(dbi)

	cursor, err := txn.OpenCursor(dbi)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close()

	limit := uint64(1000)
	for k, _, err := cursor.Get(nil, nil, mdbx.First); err == nil; k, _, err = cursor.Get(nil, nil, mdbx.Next) {
		if uint64(len(k)*2) == length {
			count++
			if count <= limit {
				keyHex := hex.EncodeToString(k)
				keys = append(keys, keyHex)
			}
		}
	}

	if err != nil && !errors.Is(err, mdbx.NotFound) {
		return 0, nil, err
	}

	return count, keys, nil
}

func (m *MDBX) CountValuesOfLength(bucketName string, length uint64) (uint64, []string, error) {
	var count uint64
	var values []string

	txn, err := m.env.BeginTxn(nil, mdbx.Readonly)
	if err != nil {
		return 0, nil, err
	}
	defer txn.Abort()

	dbi, err := txn.OpenDBISimple(bucketName, 0)
	if err != nil {
		return 0, nil, err
	}
	defer m.env.CloseDBI(dbi)

	cursor, err := txn.OpenCursor(dbi)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close()

	limit := uint64(1000)
	for _, v, err := cursor.Get(nil, nil, mdbx.First); err == nil; _, v, err = cursor.Get(nil, nil, mdbx.Next) {
		if uint64(len(v)*2) == length {
			count++
			if count <= limit {
				keyHex := hex.EncodeToString(v)
				values = append(values, keyHex)
			}
		}
	}

	if err != nil && !errors.Is(err, mdbx.NotFound) {
		return 0, nil, err
	}

	return count, values, nil
}
