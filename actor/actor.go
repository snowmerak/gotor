package actor

import (
	"bytes"
	"go/format"
	"strings"
)

func Generate(packageName string, actorName string, channelParameters map[string]string) []byte {
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString("package " + packageName + "\n\n")
	buffer.WriteString("import \"context\"\n\n")
	buffer.WriteString("type " + actorName + " struct {\n")
	buffer.WriteString("\tctx context.Context\n")
	buffer.WriteString("\tcancel context.CancelFunc\n")
	for channelName, channelType := range channelParameters {
		buffer.WriteString("\t" + channelName + " chan " + channelType + "\n")
	}
	buffer.WriteString("\t// TODO: Write your actor states here\n")
	buffer.WriteString("}\n\n")

	newFunctionName := actorName
	if strings.ToLower(packageName) == strings.ToLower(actorName) {
		newFunctionName = ""
	}

	buffer.WriteString("func New" + newFunctionName + "(ctx context.Context, queueSize int) *" + actorName + " {\n")
	buffer.WriteString("\tctx, cancel := context.WithCancel(ctx)\n")
	for channelName, channelType := range channelParameters {
		buffer.WriteString("\t" + channelName + " := make(chan " + channelType + ", queueSize)\n")
	}
	buffer.WriteString("\n")

	buffer.WriteString("\tgo func() {\n")
	buffer.WriteString("\t\tfor {\n")
	buffer.WriteString("\t\t\tselect {\n")
	buffer.WriteString("\t\t\tcase <-ctx.Done():\n")
	buffer.WriteString("\t\t\t// TODO: Write your actor stop logic here\n")
	buffer.WriteString("\t\t\treturn\n")
	for channelName := range channelParameters {
		buffer.WriteString("\t\t\tcase value := <-" + channelName + ":\n")
		buffer.WriteString("\t\t\t// TODO: Write your actor logic here\n")
	}
	buffer.WriteString("\t\t\t}\n")
	buffer.WriteString("\t\t}\n")
	buffer.WriteString("\t}()\n")

	buffer.WriteString("\treturn &" + actorName + "{\n")
	buffer.WriteString("\t\tctx: ctx,\n")
	buffer.WriteString("\t\tcancel: cancel,\n")
	for channelName := range channelParameters {
		buffer.WriteString("\t\t" + channelName + ": " + channelName + ",\n")
	}
	buffer.WriteString("\t}\n")
	buffer.WriteString("}\n\n")

	for channelName, channelType := range channelParameters {
		functionName := strings.ToUpper(channelName[:1]) + channelName[1:]
		receiverName := strings.ToLower(actorName[:1])
		buffer.WriteString("func (" + receiverName + " *" + actorName + ") " + functionName + "(value " + channelType + ") {\n")
		buffer.WriteString("\t" + receiverName + "." + channelName + "<- value\n")
		buffer.WriteString("}\n\n")
	}

	buffer.WriteString("func (m *" + actorName + ") Stop() {\n")
	buffer.WriteString("\tm.cancel()\n")
	buffer.WriteString("}\n\n")

	for channelName, channelType := range channelParameters {
		functionName := "GetSenderOf" + strings.ToUpper(channelName[:1]) + channelName[1:]
		receiverName := strings.ToLower(actorName[:1])
		buffer.WriteString("func (" + receiverName + " *" + actorName + ") " + functionName + "() chan<- " + channelType + " {\n")
		buffer.WriteString("\treturn " + receiverName + "." + channelName + "\n")
		buffer.WriteString("}\n\n")
	}

	generated := buffer.Bytes()
	result, err := format.Source(generated)
	if err != nil {
		return generated
	}

	return result
}
