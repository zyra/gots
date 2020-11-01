# GoTS

Go library and cross-platform CLI to transpile Go definitions into TypeScript types.

Useful if you are writing API services using Go for web/mobile apps built with TypeScript.

## Install CLI
> TODO: add docs

## CLI Usage
Generating code with the CLI can be configured with flags, or a yaml file. 

#### CLI Flags

#### YAML Configuration
Yaml configuration is recommended if you have a large project and you want to specify a large list of directories to include/exclude.



## Programmatic Usage
> TODO: add docs

## Concepts
GoTS uses the standard `go/ast` package to parse the provided source files and export type aliases, structs and interfaces as TypeScript code.

#### Struct output example

```go
// Basic go struct with a few properties
type Account struct {
    // basic string property that will be exported 
    // with the name defined in the json tag
    Name     string     `json:"name"`
   
    // omitempty will mark this property as optional
    Email    string     `json:"email,omitempty"`

    // this property doesn't have a json tag
    // so it will be exported with the go name
    // similar to how json marshaller works
    Friends []string
    
    // a map translates to a generic typescript object
    Labels  map[string]string `json:"labels,omitempty"`
    
    // interface{} translates to any
    Metadata    map[string]interface{} `json:"metadata"`

    // time will be exported as a string to match the json output
    UpdatedAt   time.Time   `json:"updatedAt"`
}
```

GoTS will generate the following TypeScript code
```ts
export interface Account {
    name: string;
    email?: string;
    Friends: string[];
    labels?: { [key: string]: string };
    metadata: { [key: string]: any };
    updatedAt: string;
}
```

#### Struct type override

In some cases your struct may contain an external type, unexported type, 
or a type that transforms to a different JSON type.
You can use `gots` tags in your struct definition to specify the desired output type.

Example:
```go
import "go.mongodb.org/mongo-driver/bson/primitive"

type Image struct {
    ID  primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty" gots:"type:string"`
    URL string              `json:"url"`
}
```

Output:
```ts
export interface Image {
    id?: string;
    url: string;
}
```

## Output styles
TypeScript projects are built in different ways and this library tries to provide different ways of exporting 
TypeScript code to work with existing project structures.

The examples below will use the following go package structure as the input:
```
account/        # account package
    account.go  # contains Account struct

image/              # image package
    image.go        # contains Image struct
    image_type.go   # contains ImageType type alias
```

#### All-in-one file
Exports all parsed types to a single `.ts` file.
This mode might not work in all cases if you are generating types recursively and your packages contain duplicate names.

Generated file:
```
models.ts    # contains all parsed types
```
Usage example:
```ts
import * as models from './models'
let account: models.Account;
let image: models.Image;
let imageType: models.ImageType;

// or
import { Account, Image, ImageType } from './models';
let account: Account;
let image: Image;
let imageType: ImageType;
```

#### File per package
Produces a separate file for each go package found.
This is the default behaviour of the CLI.

Generated files:
```
account.ts  # contains Account interface
image.ts    # contains Image interface and ImageType type alias
```

Usage example:
```ts
import { Account } from './account';
import { Image, ImageType } from './image';

let account: Account;
let image: Image;
let imageType: ImageType;

// or if you prefer to refer to types in a way that matches your go code
import * as account from './account';
import * as image from './image';

let account: account.Account;
let image: image.Image;
let imageType: image.ImageType;
```

#### Mirror go package
Produces a directory structure similar to your go package, with individual files that match your go package structure.

Generated files:
```
account/
    account.ts    # contains Account interface
    index.ts      # barrel export
image/
    image.ts      # contains Image interface
    image_type.ts # contains ImageType type alias
    index.ts      # barrel export
```

Usage example:
```ts
import { Account } from './account/account';
import { Image } from './image/image';
import { ImageType } from './image/image_type';

let account: Account;
let image: Image;
let imageType: ImageType;

// or if you prefer to refer to types in a way that matches your go code
import * as account from './account';
import * as image from './image';

let account: account.Account;
let image: image.Image;
let imageType: image.ImageType;
```