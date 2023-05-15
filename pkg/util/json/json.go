package json

import (
	"fmt"
	"reflect"
)

// func GetFromYaml(filename string, a interface{}) error {
// 	// 读取Pod YAML文件
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()
// 	// Read the YAML file
// 	yamlFile, err := io.ReadAll(file)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println("The file is:\n", string(yamlFile))

// 	// 解析YAML文件
// 	var yamlObj unstructured.Unstructured
// 	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(yamlFile), len(yamlFile))
// 	if err := decoder.Decode(&yamlObj); err != nil {
// 		panic(err)
// 	}

// 	// 将解析后的对象转换为对应的类型
// 	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(yamlObj.Object, a); err != nil {
// 		panic(err)
// 	}

// 	return nil
// }

func compareStruct(a, b interface{}) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	if av.Type() != bv.Type() {
		fmt.Println("The two structs are not the same type.")
		return
	}

	for i := 0; i < av.NumField(); i++ {
		af := av.Field(i)
		bf := bv.Field(i)

		if !reflect.DeepEqual(af.Interface(), bf.Interface()) {
			fmt.Printf("The field %q is different: %v != %v\n", av.Type().Field(i).Name, af.Interface(), bf.Interface())
		}
	}
}

func CheckDeepEqual(a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		fmt.Println("The two structs are not equal.")
		compareStruct(a, b)
	} else {
		fmt.Println("The two structs are equal.")
	}
}
