package unique

import (
	"context"
	"errors"
	"fmt"
	"site/utils/consts"
	"site/utils/log"

	"cloud.google.com/go/datastore"
)

var (
	ErrEntityAlreadyExists = errors.New("datastore: entity already exist")
)

type Constraint struct {
	Ref      string
	Inactive bool
	Kind     string `datastore:"-"`
	Value    string `datastore:"-"`
}

func (c *Constraint) UniqueKind() string {
	return fmt.Sprintf("_Unique_%s", c.Kind)
}

func (c *Constraint) RefKey() *datastore.Key {
	key, err := datastore.DecodeKey(c.Ref)
	if err != nil {
		return nil
	}
	return key
}

func Put(c context.Context, key *datastore.Key, src interface{}, constraints ...Constraint) (*datastore.Key, error) {

	dsClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Falha ao criar client do Datastore: %v", err)
		return nil, err
	}
	defer dsClient.Close()

	for i := range constraints {
		cons := constraints[i]
		log.Infof(c, "Verify if constraint exist: %#v", cons)
		if cons.Value == "" || cons.Kind == "" {
			return nil, fmt.Errorf(`Kind '%s' and value '%s' are required to create constraint`, cons.Kind, cons.Value)
		}
		var found Constraint
		uniqueKey := datastore.NameKey(cons.UniqueKind(), cons.Value, nil)
		err := dsClient.Get(c, uniqueKey, &found)
		switch err {
		case datastore.ErrNoSuchEntity:
			log.Infof(c, "No constraint found")
			continue
		case nil:
			if found.Inactive {
				log.Infof(c, "Found inactive contraint %#v", found)
				foundKey, err := datastore.DecodeKey(found.Ref)
				if err != nil {
					return nil, err
				}
				if foundKey.Kind != key.Kind {
					return nil, fmt.Errorf("Incompatible kind %s<>%s", foundKey.Kind, key.Kind)
				}
				log.Infof(c, "Updating key %s with %s", key.String(), foundKey.String())
				key = foundKey
			} else {
				if key.Encode() == found.Ref {
					log.Infof(c, "Constraint %#v found", found)
					continue
				} else {
					log.Infof(c, "Constraint ref no matches %s<>%s", found.Ref, cons.Ref)
					return nil, ErrEntityAlreadyExists
				}
			}
		default:
			log.Warningf(c, "Error on get constraint: %#v", err)
			return nil, err
		}
	}
	key, err = dsClient.Put(c, key, src)
	if err != nil {
		return nil, err
	}

	for i := range constraints {
		constraints[i].Ref = key.Encode()
		k := datastore.NameKey(constraints[i].UniqueKind(), constraints[i].Value, nil)
		log.Infof(c, "Insert constraint %v", constraints[i])
		if _, err := dsClient.Put(c, k, &constraints[i]); err != nil {
			return nil, err
		}
	}
	return key, nil
}

func Get(c context.Context, constraint *Constraint, dst interface{}) error {

	dsClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Falha ao criar client do Datastore: %v", err)
		return err
	}
	defer dsClient.Close()

	if constraint == nil {
		return fmt.Errorf("Constraint cannot be nil")
	}
	uniqueKey := datastore.NameKey(constraint.UniqueKind(), constraint.Value, nil)
	if err := dsClient.Get(c, uniqueKey, constraint); err != nil {
		return err
	}
	if dst == nil {
		return nil
	}
	key, err := datastore.DecodeKey(constraint.Ref)
	if err != nil {
		return err
	}
	return dsClient.Get(c, key, dst)
}
