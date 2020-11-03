package reader

// contains interfaces, types, and constants that were found in the same file
type File struct {
	Constants  []*Constant  `json:"constants"`  // constants
	Interfaces []*Interface `json:"interfaces"` // interfaces
	Types      []*TypeAlias `json:"types"`      // types
	SourcePath string       `json:"-"`          // original source path
}

func (f *File) EachInterface(fx func(iface *Interface) (bool, error)) error {
	for i := range f.Interfaces {
		if ok, err := fx(f.Interfaces[i]); err != nil {
			return err
		} else if !ok {
			return nil
		}
	}
	return nil
}

func (f *File) EachTypeAlias(fx func(t *TypeAlias) (bool, error)) error {
	for i := range f.Types {
		if ok, err := fx(f.Types[i]); err != nil {
			return err
		} else if !ok {
			return nil
		}
	}
	return nil
}

func (f *File) EachConst(fx func(t *Constant) (bool, error)) error {
	for i := range f.Constants {
		if ok, err := fx(f.Constants[i]); err != nil {
			return err
		} else if !ok {
			return nil
		}
	}
	return nil
}