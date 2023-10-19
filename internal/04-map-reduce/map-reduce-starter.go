package main

// play https://go.dev/play/p/0HQJ81CDwwB

var lines = []string{
	"lorem ipsum is simply dummy text of the printing and typesetting industry",
	"lorem ipsum has been the industry's standard dummy text ever since the",
	"when an unknown printer took a galley of type and scrambled it to make a type specimen book",
	"it has survived not only five centuries but also the leap into electronic typesetting remaining essentially unchanged",
	"it was popularised in the with the release of letraset sheets containing lorem ipsum passages and more recently with desktop publishing software like aldus pagemaker including versions of lorem ipsum",
}

func main() {

	// 1. pass each line into lineChan
	// 2. start mappers which read from lineChan, split lines into words, and send via wordChannels to the appropriate
	//    reducers.
	// 3. start reducers which read from wordChannels, form a local wordcount, and send the results to countChannel.
	//
	// Make sure everything shuts down gracefully before the process ends.
}
