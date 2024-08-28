package main

import (
	"fmt"
	"go/ast"
)

func (s *Source) GenerateRemoteImpl() error {
	data, err := GenerateRemoteImpl(s.FQDN(), s.RemoteName(), s.astFile, s.spec)
	if err != nil {
		return err
	}

	var b = &s.file.common
	b.Sprintf("// Generated from %s\n", s.fileName)
	b.WriteString(data)
	b.WriteString("\n")
	return nil
}

/*

// Generated from room_service.go
type RoomServiceRemoteImpl struct{}
type RoomService_AddUser_Request struct {
	RoomId string `bson:"roomId"`
	UserId string `bson:"userId"`
}

func (r *RoomServiceRemoteImpl) AddUser(ctx context.Context, roomId string, userId string) (bool, error) {
	data, err := bson.Marshal(RoomService_AddUser_Request{roomId, userId})
	if err != nil {
		panic(err)
	}
	request := veil.Request{
		Service: "main.RoomService",
		Method:  "AddUser",
		Args:    data,
	}
	reply := []any{}
	var result0 bool
	var result1 error
	err = veil.Call(request, &reply)
	if err != nil {
		result1 = err
	} else {
		result0 = veil.NilGet[bool](reply[0])
		result1 = veil.NilGet[error](reply[1])
	}
	return result0, result1
}

type RoomService_RemoveUser_Request struct {
	RoomId string `bson:"roomId"`
	UserId string `bson:"userId"`
}

func (r *RoomServiceRemoteImpl) RemoveUser(ctx context.Context, roomId string, userId string) (bool, error) {
	data, err := bson.Marshal(RoomService_RemoveUser_Request{roomId, userId})
	if err != nil {
		panic(err)
	}
	request := veil.Request{
		Service: "main.RoomService",
		Method:  "RemoveUser",
		Args:    data,
	}
	reply := []any{}
	var result0 bool
	var result1 error
	err = veil.Call(request, &reply)
	if err != nil {
		result1 = err
	} else {
		result0 = veil.NilGet[bool](reply[0])
		result1 = veil.NilGet[error](reply[1])
	}
	return result0, result1
}

type RoomService_Broadcast_Request struct {
	RoomId string `bson:"roomId"`
	Msg    string `bson:"msg"`
}

func (r *RoomServiceRemoteImpl) Broadcast(ctx context.Context, roomId string, msg string) (bool, error) {
	data, err := bson.Marshal(RoomService_Broadcast_Request{roomId, msg})
	if err != nil {
		panic(err)
	}
	request := veil.Request{
		Service: "main.RoomService",
		Method:  "Broadcast",
		Args:    data,
	}
	reply := []any{}
	var result0 bool
	var result1 error
	err = veil.Call(request, &reply)
	if err != nil {
		result1 = err
	} else {
		result0 = veil.NilGet[bool](reply[0])
		result1 = veil.NilGet[error](reply[1])
	}
	return result0, result1
}


*/

func generateRequestStruct(typeSpec *ast.TypeSpec, method *ast.FuncDecl) (string, string) {
	var b Builder

	// Create requests object
	reqObjName := fmt.Sprintf("%s_%s_Request", typeSpec.Name.Name, method.Name)
	b.Sprintf("type %s struct {\n", reqObjName)
	for i, param := range method.Type.Params.List {
		for _, name := range param.Names {
			if i != 0 {
				tas := getTypeAsString(param.Type)
				b.Sprintf("%s %s\n", UppercaseFirst(name.Name), tas)
			}
		}
	}
	b.WriteString("}\n")

	return b.String(), reqObjName
}

func generateSerialization(requestObjName string, method *ast.FuncDecl) string {
	var b Builder

	temp := ""
	for i, param := range method.Type.Params.List {
		for _, name := range param.Names {
			if i != 0 {
				if i > 1 {
					temp += ","
				}
				temp += name.Name
			}
		}
	}
	b.Sprintf("data, err := bson.Marshal(%s{ %s})\n", requestObjName, temp)
	b.Sprintf("if err != nil {\n")
	b.Sprintf("	panic(err)\n")
	b.Sprintf("}\n")

	return b.String()
}

func GenerateRemoteImpl(fqdn string, remoteImplName string, file *ast.File, typeSpec *ast.TypeSpec) (string, error) {
	// Generate interface name based on the struct name.

	var b Builder

	b.Sprintf("type %s struct {}\n", remoteImplName)

	// Iterate over the methods and generate method signatures.
	for _, method := range GetMethodsForStruct(file, typeSpec.Name.Name) {
		methodSignature := GenerateMethodSignature(method)
		if methodSignature != "" {

			requestDecl, requestObjName := generateRequestStruct(typeSpec, method)
			b.WriteString(requestDecl)

			b.Sprintf("func (r *%s) %s {\n", remoteImplName, methodSignature)
			b.WriteString(generateSerialization(requestObjName, method))

			b.Sprintf("request := veil.Request{\n")
			b.Sprintf("	Service: \"%s\",\n", fqdn)
			b.Sprintf("	Method:  \"%s\",\n", method.Name.String())
			b.Sprintf("	Args:    data,\n")
			b.Sprintf("}\n")
			b.Sprintf("reply := []any{}\n")

			// Handle the return values.
			last := ""
			if method.Type.Results != nil {
				for i, result := range method.Type.Results.List {
					tas := getTypeAsString(result.Type)
					b.Sprintf("var result%d %s\n", i, tas)
					last = fmt.Sprintf("result%d", i)
				}
			}
			b.Sprintf("err = veil.Call(request, &reply)\n")

			b.Sprintf("if err != nil {\n")
			b.Sprintf("	%s = err\n", last)
			b.Sprintf("} else {\n")

			if method.Type.Results != nil {
				for i, result := range method.Type.Results.List {
					tas := getTypeAsString(result.Type)
					b.Sprintf("result%d = veil.NilGet[%s](reply[%d])\n", i, tas, i)
				}
			}
			b.Sprintf("}\n")

			b.Sprintf("return ")
			if method.Type.Results != nil {
				for i := range method.Type.Results.List {
					if i == 0 {
						b.Sprintf("result%d", i)
					} else {
						b.Sprintf(", result%d", i)

					}
				}
			}
			b.Sprintf("\n")

			b.Sprintf("}\n")
		}
	}

	return b.String(), nil
}
