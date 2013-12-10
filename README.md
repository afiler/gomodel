gomodel
=======
gomodel is a library for easy mapping between SQL databases and Go structs.

### Defining a record structure:
```
type Trademark struct {
	Agency string
	Serial int
	Mark string
	TrademarkClasses string
	Record
}
```
Record is provided by gomodel for semi-easy mixins of convenience methods, and
for identifying structs that are meant to be used with gomodel's Model class.

### Defining a model for the record:
```
func TrademarkModel(config Config) Model {
	return Model {
		Data: &Trademark{},
		Prototype: Trademark{},
		TableName: "trademarks",
		Config: config,
	}
}
```
The model keeps configuration parameters and holds a prototype struct to be
copied when constructing records retrieved from the database.

### Mixing in String() capability from Record
```
func (tm Trademark) String() string {
	return tm.StringObj(&m)
}
```

### Querying
```
config := Config { DriverName: "somedriver", DataSourceName: "someuser:somepassword@/somedb" }

trademarkModel := MarkModel(config)
tm := trademarkModel.Find(struct {
	Agency string
	Serial int
}{"US", 70011210})

fmt.Println(tm)
```
