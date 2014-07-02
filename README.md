jsonproperty
============

A datastore JsonProperty in Go for Google App Engine.  Equivalent to the Python
[ndb.JsonProperty](https://developers.google.com/appengine/docs/python/ndb/properties).

Usage
-----

This package contains a JsonProperty for datastore structs that should be used in conjuction with
the struct's [Load() and Save() methods](https://developers.google.com/appengine/docs/go/datastore/reference#PropertyLoadSaver).

```go
import (
        "appengine/datastore"
        "github.com/matthewbelisle-wf/jsonproperty"
        "time"
)

type Person struct {
        Name     string                    `datastore:"name"`
        Birth    time.Time                 `datastore:"birth"`
        MetaData jsonproperty.JsonProperty `datastore:"-" jsonproperty:"metadata"`
}

func (person *Person) Load(c <-chan datastore.Property) error {
        c, err := jsonproperty.LoadJsonProperties(person, c)
        if err != nil {
                return err
        }
        return datastore.LoadStruct(person, c)
}

func (person *Person) Save(c chan<- datastore.Property) error {
        c, err := jsonproperty.SaveJsonProperties(person, c)
        if err != nil {
                return err
        }
        return datastore.SaveStruct(person, c)
}
```

It saves the property as a `string` in the datastore but converts it to a `map[string]interface{}`
on load.  In order to work, the normal datastore `LoadStruct()` and `SaveStruct()` have to know to
ignore it via the datastore tag `datastore:"-"`.  For more details see the
[example](./example/app.go).