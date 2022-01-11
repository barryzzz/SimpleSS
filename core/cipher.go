package core

type Cipher struct {
	encodePassword *Password
	decodePassword *Password
}

func (cipher *Cipher) encode(bs []byte) {
	for i, v := range bs {
		bs[i] = cipher.encodePassword[v]
	}
}

func (cipher *Cipher) decode(bs []byte) {
	for i, v := range bs {
		bs[i] = cipher.decodePassword[v]
	}
}

func NewCipher(pwd *Password) *Cipher {
	dec := &Password{}
	enc := &Password{}
	for i, v := range pwd {
		dec[i] = v
		enc[v] = byte(i)
	}
	return &Cipher{enc, dec}
}

func (cipher *Cipher) GetEncode() *Password {
	return cipher.encodePassword
}

func (cipher *Cipher) GetDecode() *Password {
	return cipher.decodePassword
}
