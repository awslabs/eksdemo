package key

import (
	"fmt"
	"os"
	"sort"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type KMSKey struct {
	Aliases []types.AliasListEntry
	Key     *types.KeyMetadata
}

type Getter struct {
	kmsClient *aws.KMSClient
}

func NewGetter(kmsClient *aws.KMSClient) *Getter {
	return &Getter{kmsClient}
}

func (g *Getter) Init() {
	if g.kmsClient == nil {
		g.kmsClient = aws.NewKMSClient()
	}
}

func (g *Getter) Get(alias string, output printer.Output, options resource.Options) error {
	kmsOptions, ok := options.(*Options)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to KmsKeyOptions")
	}
	_ = kmsOptions

	var keys []*KMSKey
	var key *KMSKey
	var err error

	if alias != "" {
		key, err = g.GetByAlias(alias)
		keys = []*KMSKey{key}
	} else {
		keys, err = g.GetAllKeys()
	}
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(keys))
}

func (g *Getter) GetAllKeys() ([]*KMSKey, error) {
	keyList, err := g.kmsClient.ListKeys()
	if err != nil {
		return nil, err
	}

	keyAliasMapping := map[string]*KMSKey{}
	for _, k := range keyList {
		key, err := g.kmsClient.DescribeKey(awssdk.ToString(k.KeyId))
		if err != nil {
			return nil, err
		}
		keyAliasMapping[awssdk.ToString(k.KeyId)] = &KMSKey{[]types.AliasListEntry{}, key}
	}

	aliases, err := g.kmsClient.ListAliases()
	if err != nil {
		return nil, err
	}

	for _, a := range aliases {
		if k, ok := keyAliasMapping[awssdk.ToString(a.TargetKeyId)]; ok {
			k.Aliases = append(k.Aliases, a)
		}
	}

	keys := make([]*KMSKey, 0, len(keyAliasMapping))
	for _, key := range keyAliasMapping {
		keys = append(keys, key)
	}

	// Show recently created Keys at the end of the list
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Key.CreationDate.Before(awssdk.ToTime(keys[j].Key.CreationDate))
	})

	return keys, nil
}

func (g *Getter) GetByAlias(aliasName string) (*KMSKey, error) {
	aliases, err := g.kmsClient.ListAliases()
	if err != nil {
		return nil, err
	}

	for _, a := range aliases {
		if "alias/"+aliasName != awssdk.ToString(a.AliasName) {
			continue
		}

		keyId := awssdk.ToString(a.TargetKeyId)

		key, err := g.kmsClient.DescribeKey(keyId)
		if err != nil {
			return nil, err
		}

		return &KMSKey{filterAliasesByKeyId(aliases, keyId), key}, nil
	}

	return nil, &resource.NotFoundByError{Type: "kms-key", Name: "alias", Value: aliasName}
}

func filterAliasesByKeyId(aliases []types.AliasListEntry, id string) []types.AliasListEntry {
	filtered := []types.AliasListEntry{}
	for _, a := range aliases {
		if awssdk.ToString(a.TargetKeyId) == id {
			filtered = append(filtered, a)
		}
	}
	return filtered
}
