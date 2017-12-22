package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	//"strings"

	yaml "gopkg.in/yaml.v2"
)

type Role string

type OperationConfig struct {
	Enable bool `json:"enable" yaml:"enable"`
}

type Operation interface {
	Name() string
	Role() Role
	Go(w http.ResponseWriter, r *http.Request) error
}

var allOps = map[Role][]Operation{}

func GetOperationNames(role Role) []string {
	roleOps := allOps[role]
	out := make([]string, len(roleOps))
	for i, op := range roleOps {
		out[i] = op.Name()
	}
	return out
}

func GetRoleNames() []string {
	out := make([]string, 0, len(allOps))
	for role := range allOps {
		out = append(out, string(role))
	}
	return out
}

func GetOperation(role Role, name string) (Operation, error) {
	roleOps, ok := allOps[role]
	if !ok {
		return nil, fmt.Errorf("Role Not Found: %s", string(role))
	}
	for _, op := range roleOps {
		if op.Name() == name {
			return op, nil
		}
	}
	return nil, fmt.Errorf("Operation Not Found: %s", name)
}

type NewOpFunc func(Role, string) (Operation, error)

var registeredOpFunc = map[string]NewOpFunc{}

func RegisterNewOp(name string, f NewOpFunc) {
	if _, ok := registeredOpFunc[name]; ok {
		fmt.Printf("Operation registered twice, second on wins: %s\n", name)
	}
	registeredOpFunc[name] = f
}

func getDirsInPath(cfgDir string) ([]string, error) {
	files, err := ioutil.ReadDir(cfgDir)
	if err != nil {
		return nil, err
	}

	dirs := []string{}
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		}
	}

	return dirs, nil
}

func InitializeOperations(cfgDir string) error {
	// Walk the cfgDir looking for roles defined
	roles, err := getDirsInPath(cfgDir)
	if err != nil {
		return err
	}

	// Walk each role looking for ops defined
	for _, role := range roles {
		roleDir := filepath.Join(cfgDir, role)
		ops, err := getDirsInPath(roleDir)
		if err != nil {
			return err
		}
		// Walk each op calling the op's 'NewOpFunc"
		for _, opName := range ops {
			opCfgDir := filepath.Join(roleDir, opName)

			opCfgFile := filepath.Join(opCfgDir, "config")
			data, err := ioutil.ReadFile(opCfgFile)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Printf("Role: %s Operation: %s : Config file not present: %s\n", role, opName, opCfgFile)
					continue
				}
				return err
			}
			oc := OperationConfig{}
			if err := yaml.Unmarshal(data, &oc); err != nil {
				return err
			}

			if !oc.Enable {
				fmt.Printf("Role: %s Operation: %s : Disabled by config\n", role, opName)
				continue
			}

			fmt.Printf("Role: %s Operation: %s : Registering\n", role, opName)
			newOpFunc, ok := registeredOpFunc[opName]
			if !ok {
				fmt.Printf("Role: %s Operation: %s : Found config but operation does not exist: %s\n", role, opName, opCfgFile)
				continue
			}

			r := Role(role)
			op, err := newOpFunc(r, opCfgDir)
			if err != nil {
				return err
			}

			allOps[r] = append(allOps[r], op)
		}
	}
	return nil
}
