// IoTDet (c) 2017 Rafal Zajac <rzajac@gmail.com> All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hq

import (
    "crypto/aes"
    "crypto/cipher"
)

// AES represents AES encoder / decoder.
type CipherAES struct {
    key    []byte // AES encryption key.
    vector []byte // AES CBC initialization vector.
    keyLen int    // AES key length.
}

// NewCipherAES returns new AES encoder / decoder.
func NewCipherAES(key, vector []byte) *CipherAES {
    return &CipherAES{
        key:    key,
        vector: vector,
        keyLen: len(key),
    }
}

func (ad *CipherAES) Encrypt(data []byte) ([]byte, error) {
    var err error

    dataLength := len(data)
    data = padRight(data, dataLength+(ad.keyLen-dataLength%ad.keyLen))

    var block cipher.Block
    if block, err = aes.NewCipher(ad.key); err != nil {
        return nil, err
    }

    mode := cipher.NewCBCEncrypter(block, ad.vector)
    mode.CryptBlocks(data, data)

    return data, nil
}

func (ad *CipherAES) Decrypt(data []byte) ([]byte, error) {
    var err error

    var block cipher.Block
    if block, err = aes.NewCipher(ad.key); err != nil {
        return nil, err
    }

    cipher.NewCBCDecrypter(block, ad.vector).CryptBlocks(data, data)

    return data, err
}

// padRight right pads byte array with zeros.
func padRight(msg []byte, length int) []byte {
    for {
        if len(msg) == length {
            return msg
        }
        msg = append(msg, 0)
    }
}
