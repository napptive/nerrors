# nerrors
Error library for Napptive projects. 

This library aims to extend the generic golang error library in order to contribute more information to the error.
* __Stack trace__ related to where the error happened in the code base
* __Type error__ compatible with GRPC
* __Parent__ with information about the parent errors

## Example of usage

- Creating a new error
```
// extended error
err := NewNotFoundError("unable to find the record")
fmt.Println(err.StackTraceToString)
```
- For a complex wrapping example:
```
// classic golang error
common := status.Errorf(codes.Aborted, "new error")
// extended error
err := NewInternalErrorFrom(common,"internal error")
fmt.Println(err.StackTraceToString)
```
## Integration with Code Climate

This project is integrated with codeclimate

[![Maintainability](https://api.codeclimate.com/v1/badges/d426ab46dd6c71fcb93b/maintainability)](https://codeclimate.com/repos/5fc8adcdd753d801b6015bf5/maintainability) [![Test Coverage](https://api.codeclimate.com/v1/badges/d426ab46dd6c71fcb93b/test_coverage)](https://codeclimate.com/repos/5fc8adcdd753d801b6015bf5/test_coverage)

## Integration with Github Actions

This project is integrated with GitHub 

![Check changes in the Main branch](https://github.com/napptive/nerrors/workflows/Check%20changes%20in%20the%20Main%20branch/badge.svg)

## License

 Copyright 2020 Napptive

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.