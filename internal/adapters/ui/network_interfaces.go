// Copyright 2025.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ui

import (
    "context"
    "net"
    "sort"
    "time"
    "fmt"
)

func GetNetworkInterfacesWithContext(ctx context.Context) ([]string, error) {
    type result struct {
        names []string
        err   error
    }
    
    ch := make(chan result, 1)
    
    go func() {
        interfaces, err := net.Interfaces()
        if err != nil {
            ch <- result{nil, fmt.Errorf("failed to get network interfaces: %w", err)}
            return
        }

        var names []string
        for _, iface := range interfaces {
            if iface.Flags&net.FlagUp != 0 {
                names = append(names, iface.Name)
            }
        }
        sort.Strings(names)
        ch <- result{names, nil}
    }()
    
    select {
    case <-ctx.Done():
        return nil, fmt.Errorf("operation timed out retrieving network interfaces")
    case res := <-ch:
        return res.names, res.err
    }
}
