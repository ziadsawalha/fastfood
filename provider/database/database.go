package database

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/GeertJohan/go.rice"
	"github.com/jarosser06/fastfood/provider"
	"github.com/jarosser06/fastfood/provider/cookbook"
	"github.com/jarosser06/fastfood/util"
)

const (
	defaultCredsdatabag = "mysql_credentials"
	defaultType         = "mysql"
	defaultRole         = "master"
)

type Database struct {
	*provider.Helpers
	Cookbook cookbook.Cookbook
	Databag  string `json:"databag,omitempty"`
	Database string `json:"database,omitempty"`
	Name     string `json:"name,omitempty"`
	Openfor  string `json:"openfor,omitempty"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
	Type     string `json:"type,omitempty"`
	User     string `json:"user,omitempty"`
}

func New(name string, ckbk cookbook.Cookbook) Database {

	return Database{
		Cookbook: ckbk,
		Name:     name,
		Role:     defaultRole,
		Type:     defaultType,
	}
}

func (d *Database) SetOrReturnDatabase(str string) string {
	if str == "" {
		return d.QString(d.Database)
	} else {
		return d.QString(str)
	}
}

func (d *Database) Dependencies() []string {
	deps := []string{"rackspace_iptables"}
	switch d.Type {
	case "mysql":
		deps = append(deps, "mysql-multi")
	case "postgres":
		deps = append(deps, "postgres-multi")

		if d.Openfor != "" {
			deps = append(deps, "chef-sugar")
		}
	}

	if d.Database != "" {
		deps = append(deps, "database")
	}

	return deps
}

func (d *Database) GenFiles() error {
	recipeFile := fmt.Sprintf("recipes/%s.rb", d.Name)
	specFile := fmt.Sprintf("test/unit/spec/%s_spec.rb", d.Name)

	if util.FileExist(path.Join(d.Cookbook.Path, recipeFile)) {
		return errors.New(fmt.Sprintf("%s already exists", recipeFile))
	}

	cookbookFiles := map[string]string{
		recipeFile: fmt.Sprintf("recipes/%s_%s.rb", d.Type, d.Role),
		specFile:   fmt.Sprintf("test/unit/spec/%s_spec.rb", d.Type),
	}

	templateBox, _ := rice.FindBox("../templates/database")
	for cookbookFile, templateFile := range cookbookFiles {
		tmpStr, _ := templateBox.String(templateFile)

		t, err := provider.NewTemplate(cookbookFile, d, []string{tmpStr})

		if err != nil {
			return errors.New(fmt.Sprintf("Error creating template: %v", err))
		}

		t.CleanNewlines()

		if err := t.Flush(path.Join(d.Cookbook.Path, cookbookFile)); err != nil {
			return errors.New(fmt.Sprintf("Error writing file: %v", err))
		}
	}
	return nil
}

func (d *Database) GenDirs() error {
	dirs := [2]string{
		"recipes",
		"test/unit/spec",
	}

	for _, dir := range dirs {
		fullPath := path.Join(d.Cookbook.Path, dir)

		if !util.FileExist(fullPath) {
			err := os.MkdirAll(path.Join(d.Cookbook.Path, dir), 0755)

			if err != nil {
				return errors.New(fmt.Sprintf("database.GenDirs(): %v", err))
			}
		}
	}

	return nil
}
