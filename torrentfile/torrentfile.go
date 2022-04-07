package torrentfile

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"os"

	bencode "github.com/jackpal/bencode-go"
)

const (
	HashLength = 20
)

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

// Maybe extend to include more params
type bencodeTorrent struct {
	Announce     string      `bencode:"announce"`
	Info         bencodeInfo `bencode:"info"`
	Comment      string      `bencode:"comment"`
	CreationDate int         `bencode:"creation date"`
}

type TorrentFile struct {
	Announce    string
	InfoHash    [HashLength]byte
	PieceHashes [][HashLength]byte
	PieceLength int
	Length      int
	Name        string
}

// Computes the SHA-1 hash value of the torrent info dictionary
func (bTorrentInfo bencodeInfo) computeHash() ([HashLength]byte, error) {
	// Creates a new byte buffer to store results
	buf := new(bytes.Buffer)

	// Write the encoded value of bencodeInfo to the bytes buffer
	err := bencode.Marshal(buf, bTorrentInfo)

	if err != nil {
		return [HashLength]byte{}, err
	}

	// Computes the SHA-1 Hash
	hash := sha1.Sum(buf.Bytes())

	return hash, nil
}

// Splits the hash values stored in Info.Pieces
func (bTorrentInfo bencodeInfo) computePieceHashes() ([][HashLength]byte, error) {
	// Creates a new byte buffer to store results

	// Convert the Pieces (string) to byte array
	buf := []byte(bTorrentInfo.Pieces)

	if len(buf)%HashLength != 0 {
		err := fmt.Errorf("received invalid info.pieces of length %v", len(buf))
		return nil, err
	}

	// Compute the amount of hashes we have
	numHashes := len(buf) / HashLength
	hashes := make([][HashLength]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		// Start and End index to copy from
		start := i * HashLength
		end := start + HashLength

		copy(hashes[i][:], buf[start:end])
	}

	return hashes, nil
}

// Convers the bencoded data to a more usable struct
func (bTorrent bencodeTorrent) toTorrentFile() (TorrentFile, error) {
	// Compute the SHA-1 hash of the Info Structj
	infoHash, err := bTorrent.Info.computeHash()
	if err != nil {
		log.Panic(err)
	}

	// Compute the hashes of the Info.Pieces
	pHashes, err := bTorrent.Info.computePieceHashes()
	if err != nil {
		log.Panic(err)
	}

	torrent := TorrentFile{
		Announce:    bTorrent.Announce,
		InfoHash:    infoHash,
		PieceHashes: pHashes,
		PieceLength: bTorrent.Info.PieceLength,
		Length:      bTorrent.Info.Length,
		Name:        bTorrent.Info.Name,
	}

	return torrent, nil
}

// Opens a .torrent file and parses it
func Open(filePath string) (TorrentFile, error) {
	// Opens the file
	file, err := os.Open(filePath)

	// If any error occurs, return the error
	if err != nil {
		return TorrentFile{}, err
	}

	// Close file when done
	defer file.Close()

	// Instantiate a new empty struct to read bencode info in
	bTorrent := bencodeTorrent{}
	err = bencode.Unmarshal(file, &bTorrent)

	fmt.Println("Bencoded Torrent")

	// If any error occurs, return the error
	if err != nil {
		return TorrentFile{}, err
	}

	return bTorrent.toTorrentFile()
}
