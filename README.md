[![Documentation](https://godoc.org/github.com/HewlettPackard/simplivity-go/ovc?status.svg)](https://godoc.org/github.com/HewlettPackard/simplivity-go/ovc)
[![Build Status](https://travis-ci.com/HewlettPackard/simplivity-go.svg?branch=master)](https://travis-ci.com/HewlettPackard/simplivity-go)
[![Coverage Status](https://coveralls.io/repos/github/HewlettPackard/simplivity-go/badge.svg?branch=master)](https://coveralls.io/github/HewlettPackard/simplivity-go?branch=master)

# HPE SimpliVity Go SDK

This library provides a Go interface to the HPE SimpliVity REST APIs.

HPE SimpliVity is an intelligent hyperconverged platform that speeds application performance,
improves efficiency and resiliency, and backs up and restores VMs in seconds.

## Usage
```
import "github.com/HewlettPackard/simplivity-go/ovc"
```
Construct a new OVC client, then use the various resource clients available on the OVC client to access the resource specific features. For example:
```
//Create an OVC client
client, _ := ovc.NewClient("username", "password", "ovc_ip", "certificate_path") // with certificate
client, _ := ovc.NewClient("username", "password", "ovc_ip", "") // without certificate

//Get all the VMs without Filters
vmList, _ := client.VirtualMachines.GetAll(ovc.GetAllParams{})

//Get a VM resource by its name
vmByName, _ = client.VirtualMachines.GetByName(vmName)

//Clone the above VM
vm, err = vmByName.Clone("new_vm_name", false)
```
For more examples, head over to the [example](examples) directory.

## API Implementation

Status of the HPE SimpliVity REST interfaces that have been implemented in this Go library can be found in the [endpoints-support](endpoints-support.md) file.

## Contributing and feature requests

**Contributing:** We welcome your contributions to the Go SDK for HPE SimpliVity. See [CONTRIBUTING.md](CONTRIBUTING.md) for more details.

**Feature Requests:** If you have a need that is not met by the current implementation, please let us know (via a new issue).
This feedback is crucial for us to deliver a useful product. Do not assume that we have already thought of everything, because we assure you that is not the case.

#### Testing

When contributing code to this project, we require tests to accompany the code being delivered.
That ensures a higher standing of quality, and also helps to avoid minor mistakes and future regressions.

## License

This project is licensed under the Apache license. Please see [LICENSE](LICENSE) for more information.

## Version and changes

To view history and notes for this version, view the [CHANGELOG](CHANGELOG.md).
