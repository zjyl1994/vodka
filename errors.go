package vodka

import "fmt"

var errList = map[string]int{
	"Internal": 500,
}

type Error struct{
	identifier string
	errorMsg string
}

var NoError = Error {
	identifier : "-",
	errorMsg: "",
}

func NewMessageError(identifier,format string,a ...interface{}) Error{
	if identifier == "-"{
		return NoError
	}else{
		return Error{
			identifier:identifier,
			errorMsg:fmt.Sprintf(format,a...),
		}
	}
}

func NewError(err error) Error{
	if err == nil{
		return NoError
	}else{
		return Error{
			identifier:"Internal",
			errorMsg:err.Error(),
		}
	}
}

func (e *Error) Error()string{
	return e.errorMsg
}

func (e *Error) Identifier()string{
	return e.identifier
}

func RegisterError(identifier string,statusCode int){
	if identifier == "-"{
		Logger.Fatalln("Register reserved identifier")
	}else{
		_, exist := errList[identifier]
		if exist{
			Logger.Fatalln("Register duplicate identifier")
		}else{
			errList[identifier] = statusCode
		}
	}
}

func MultiRegisterError(errors map[string]int){
	for identifier,statusCode := range errors{
		RegisterError(identifier,statusCode)
	}
}

func hasError(err Error) bool{
	return err.identifier != "-"
}