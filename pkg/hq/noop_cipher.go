// HomeHQ (c) 2018 Rafal Zajac <rzajac@gmail.com> All rights reserved.
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

// NoopCipher represents no operation encoder decoder.
type NoopCipher struct{}

// NewNoopCipher returns new instance of CipherNoop.
func NewNoopCipher() *NoopCipher {
    return &NoopCipher{}
}

func (nd *NoopCipher) Encrypt(data []byte) ([]byte, error) {
    return data, nil
}

func (nd *NoopCipher) Decrypt(data []byte) ([]byte, error) {
    return data, nil
}
