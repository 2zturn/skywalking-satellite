// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package grpc

import (
	"fmt"

	pbv1 "github.com/golang/protobuf/proto"
	"google.golang.org/grpc/encoding"
	"google.golang.org/protobuf/proto"
)

func init() {
	encoding.RegisterCodec(codec{})
}

// OriginalData is keep binary Content
type OriginalData struct {
	Content []byte
}

func NewOriginalData(data []byte) *OriginalData {
	return &OriginalData{Content: data}
}

// codec is overwritten the original "proto" codec, and support using OriginalData to skip data en/decoding.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	original, ok := v.(*OriginalData)
	if ok {
		return original.Content, nil
	}

	vv, ok := v.(proto.Message)
	if ok {
		return proto.Marshal(vv)
	}

	b, err := pbv1.Marshal(pbv1.MessageV1(v))
	if err == nil {
		return b, nil
	}
	return nil, fmt.Errorf("failed to marshal, message is %T, want proto.MessageV1/V2 or grpc.OriginalData", v)

}

func (codec) Unmarshal(data []byte, v interface{}) error {
	original, ok := v.(*OriginalData)
	if ok {
		original.Content = data
		return nil
	}

	vv, ok := v.(proto.Message)
	if ok {
		return proto.Unmarshal(data, vv)
	}

	err := proto.Unmarshal(data, pbv1.MessageV2(v))
	if err != nil {
		return fmt.Errorf("failed to unmarshal, message is %T, want proto.MessageV1/V2 or grpc.OriginalData", v)
	}
	return nil

}

func (codec) Name() string {
	return "proto"
}
