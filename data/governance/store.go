package governance

import (
	"encoding/binary"

	"github.com/pkg/errors"

	ethchain "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

const (
	ADMIN_INITIAL_KEY    string = "initial"
	ADMIN_CURRENCY_KEY   string = "currency"
	ADMIN_FEE_OPTION_KEY string = "feeopt"

	ADMIN_EPOCH_BLOCK_INTERVAL string = "epoch"

	ADMIN_ETH_CHAINDRIVER_OPTION string = "ethcdopt"
)

type Store struct {
	state  *storage.State
	prefix []byte
}

func NewStore(prefix string, state *storage.State) *Store {
	return &Store{
		state:  state,
		prefix: storage.Prefix(prefix),
	}
}

func (st *Store) WithState(state *storage.State) *Store {
	st.state = state
	return st
}

func (st *Store) Get(key []byte) ([]byte, error) {
	prefixKey := append(st.prefix, storage.StoreKey(key)...)

	return st.state.Get(prefixKey)
}

func (st *Store) Set(key []byte, value []byte) error {
	prefixKey := append(st.prefix, storage.StoreKey(key)...)
	err := st.state.Set(prefixKey, value)
	return err
}

func (st *Store) Exists(key []byte) bool {
	prefixKey := append(st.prefix, storage.StoreKey(key)...)
	return st.state.Exists(prefixKey)
}

func (st *Store) GetCurrencies() (balance.Currencies, error) {
	result, err := st.Get([]byte(ADMIN_CURRENCY_KEY))
	currencies := make(balance.Currencies, 0, 10)
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(result, &currencies)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the currencies")
	}
	return currencies, nil
}

func (st *Store) SetCurrencies(currencies balance.Currencies) error {

	currenciesBytes, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(currencies)
	if err != nil {
		return errors.Wrap(err, "failed to serialize currencies")
	}
	err = st.Set([]byte(ADMIN_CURRENCY_KEY), currenciesBytes)
	if err != nil {
		return errors.Wrap(err, "failed to set the currencies")
	}
	return nil
}

func (st *Store) GetFeeOption() (*fees.FeeOption, error) {
	feeOpt := &fees.FeeOption{}
	bytes, err := st.Get([]byte(ADMIN_FEE_OPTION_KEY))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get FeeOption")
	}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(bytes, feeOpt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize FeeOption stored")
	}

	return feeOpt, nil
}

func (st *Store) SetFeeOption(feeOpt fees.FeeOption) error {

	bytes, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(feeOpt)
	if err != nil {
		return errors.Wrap(err, "failed to serialize FeeOption")
	}
	err = st.Set([]byte(ADMIN_FEE_OPTION_KEY), bytes)
	if err != nil {
		return errors.Wrap(err, "failed to set the FeeOption")
	}
	return nil
}

func (st *Store) Initiated() bool {
	_ = st.Set([]byte(ADMIN_INITIAL_KEY), []byte("initialed"))
	return true
}

func (st *Store) InitialChain() bool {
	data, err := st.Get([]byte(ADMIN_INITIAL_KEY))
	if err != nil {
		return true
	}
	if data == nil {
		return true
	}
	return false
}

func (st *Store) GetEpoch() (int64, error) {
	result, err := st.Get([]byte(ADMIN_EPOCH_BLOCK_INTERVAL))
	if err != nil {
		return 0, err
	}

	epoch := int64(binary.LittleEndian.Uint64(result))

	return epoch, nil
}

func (st *Store) SetEpoch(epoch int64) error {

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(epoch))

	err := st.Set([]byte(ADMIN_EPOCH_BLOCK_INTERVAL), b)
	if err != nil {
		return errors.Wrap(err, "failed to set the currencies")
	}
	return nil
}

func (st *Store) GetETHChainDriverOption() (*ethchain.ChainDriverOption, error) {
	bytes, err := st.Get([]byte(ADMIN_ETH_CHAINDRIVER_OPTION))
	if err != nil {
		return nil, err
	}
	r := &ethchain.ChainDriverOption{}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(bytes, r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize eth chaindriver option stored")
	}
	return r, nil
}

func (st *Store) SetETHChainDriverOption(opt ethchain.ChainDriverOption) error {
	bytes, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(opt)
	if err != nil {
		return errors.Wrap(err, "failed to serialize eth chaindriver option")
	}
	err = st.Set([]byte(ADMIN_ETH_CHAINDRIVER_OPTION), bytes)
	if err != nil {
		return errors.Wrap(err, "failed to set the eth chaindriver option")
	}
	return nil
}
