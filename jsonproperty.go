package jsonproperty

import (
	"appengine/datastore"
	"encoding/json"
	"reflect"
)

type JsonProperty map[string]interface{}

var jsonPropertyType = reflect.TypeOf(JsonProperty{})

// entity must be a pointer to a struct
func LoadJsonProperties(entity interface{}, c <-chan datastore.Property) (<-chan datastore.Property, error) {
	jsonProperties := map[string]reflect.Value{} // {name: value}

	// Builds jsonProperties
	value := reflect.ValueOf(entity).Elem()
	for i := 0; i < value.NumField(); i++ {
		value2 := value.Field(i)
		if value2.Type() == jsonPropertyType {
			jsonProperties[nameFromValue(value.Type().Field(i))] = value2
		}
	}

	// Builds return channel
	c2 := make(chan datastore.Property, value.NumField()-len(jsonProperties))
	defer close(c2)
	for property := range c {
		if jsonValue, ok := jsonProperties[property.Name]; ok {
			bytes := []byte(property.Value.(string))
			if err := json.Unmarshal(bytes, jsonValue.Addr().Interface()); err != nil {
				return c2, err
			}
		} else {
			c2 <- property
		}
	}

	return c2, nil
}

// entity must be a pointer to a struct
func SaveJsonProperties(entity interface{}, c chan<- datastore.Property) (chan<- datastore.Property, error) {
	value := reflect.ValueOf(entity).Elem()
	for i := 0; i < value.NumField(); i++ {
		value2 := value.Field(i)
		if value2.Type() == jsonPropertyType {
			bytes, err := json.Marshal(value2.Interface())
			if err != nil {
				return c, err
			}
			c <- datastore.Property{
				Name:  nameFromValue(value.Type().Field(i)),
				Value: string(bytes),
			}
		}
	}

	return c, nil
}

func nameFromValue(f *reflect.Value) string {
	if name := *f.Tag.Get("jsonproperty"); name != "" {
		return name
	}
	return f.Name
}
