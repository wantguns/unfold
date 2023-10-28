package ledger

import (
	"bufio"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

func GetPresentUUIDs(filename string) map[string]bool {
	file, err := os.Open(filename)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to open ledger file for parsing: %+v", filename)
	}

	defer file.Close()

	var uuids map[string]bool
	uuids = make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, "unfold_meta_start") || strings.Contains(text, "unfold_meta_end") {
			uuids[strings.Split(text, "\t")[2]] = true
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error().Err(err).Msgf("Unexpected scanner error: %+v", filename)
	}

	return uuids
}
