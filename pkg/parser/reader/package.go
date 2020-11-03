package reader

// Package represents a group of files that exist in the same file/module
type Package struct {
	Name       string  `json:"name"`  // Package name
	Files      []*File `json:"files"` // Files found in this package
	SourcePath string  `json:"-"`     // original source path
}

// Loop through all package files
func (pkg *Package) EachFile(fx func(f *File) (bool, error)) error {
	for i := range pkg.Files {
		if ok, err := fx(pkg.Files[i]); err != nil {
			return err
		} else if !ok {
			return nil
		}
	}
	return nil
}

// Loop through all package files
func (pkg *Package) EachInterface(fx func(f *File, iface *Interface) (bool, error)) error {
	return pkg.EachFile(func(f *File) (bool, error) {
		return true, f.EachInterface(func(iface *Interface) (bool, error) {
			return fx(f, iface)
		})
	})
}

// Loop through all package constants
func (pkg *Package) EachConstant(fx func(f *File, c *Constant) (bool, error)) error {
	return pkg.EachFile(func(f *File) (bool, error) {
		return true, f.EachConst(func(t *Constant) (bool, error) {
			return fx(f, t)
		})
	})
}

// Loop through all package constants
func (pkg *Package) EachTypeAlias(fx func(f *File, t *TypeAlias) (bool, error)) error {
	return pkg.EachFile(func(f *File) (bool, error) {
		return true, f.EachTypeAlias(func(t *TypeAlias) (bool, error) {
			return fx(f, t)
		})
	})
}

// Get all package constants
func (pkg *Package) Constants() []*Constant {
	ca := make([]*Constant, 0)
	_ = pkg.EachConstant(func(f *File, c *Constant) (bool, error) {
		ca = append(ca, c)
		return true, nil
	})
	return ca
}