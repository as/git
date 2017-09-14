package git

//	wire9 gitTran n[4,stringInt] data[n]
//	wire9 zobject typelen[,abomination] data[typelen]
//	wire9 packhdr sig[4,,BE] ver[4,,BE] n[4,,BE] objects[n,[]zobject]

// Object represnets a file's contents stored as a sha1 hash
type Object interface {
	Hash() Hash
	Data() []byte
}
