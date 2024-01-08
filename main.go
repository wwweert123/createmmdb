// asn-writer is an example of how to create an ASN MaxMind DB file from the
// GeoLite2 ASN CSVs. You must have the CSVs in the current working directory.
package main

import (
	"encoding/binary"
	"log"
	"net"
	"os"
	"strconv"
	"math/rand"
	"time"
	crand "crypto/rand"
	"encoding/hex"
	"flag"


	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)

func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}

func genRanIPv6part() (string, error) {
	bytes := make([]byte, 2)
  	if _, err := crand.Read(bytes); err != nil {
    	return "", err
  	}
  	return hex.EncodeToString(bytes), nil
}

func genranIPv6() string {
	s := ""
	for i := 0; i < 7; i++ {
		val, _ := genRanIPv6part()
		s += val + ":";
	}
	val, _ := genRanIPv6part()
	s += val
	return s
}

func insert_ipv4_in_sequence(n int, unique_data bool, writer *mmdbwriter.Tree) {
	for i := 16843009; i <= 16843009 + n; i++ {
		ip := int2ip(uint32(i))
		_, network, err:= net.ParseCIDR(ip.String() + "/32")
		if err != nil {
			log.Fatal(err)
		}
		
		record := mmdbtype.Map{}
		if unique_data {
			record["test"] = mmdbtype.String(strconv.Itoa(int(i)))
		} else {
			record["test"] = mmdbtype.String("hi")
		}
		
		err = writer.Insert(network, record)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func insert_ipv4_rand(n int, writer *mmdbwriter.Tree) {
	for i := 0; i <= n; i++ {
		rand.Seed(time.Now().UnixNano())
    	min := 16843009
    	max := 4294967295
	    x := rand.Intn(max - min + 1) + min
		ip := int2ip(uint32(x))
		_, network, err:= net.ParseCIDR(ip.String() + "/32")
		if err != nil {
			log.Fatal(err)
		}
		
		record := mmdbtype.Bool(true)
		// record["test"] = mmdbtype.String("hi")
		
		err = writer.Insert(network, record)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func insert_ipv6_rand(n int, writer *mmdbwriter.Tree) {
	for i := 0; i <= n; i++ {
		s := genranIPv6()
		_, network, err:= net.ParseCIDR(s + "/128")
		if err != nil {
			log.Fatal(err)
		}
		
		// record := mmdbtype.Map{}
		// record["test"] = mmdbtype.String(strconv.Itoa(int(i)))
		// record["test"] = mmdbtype.String("hi")
		record := mmdbtype.Bool(true)
		
		err = writer.Insert(network, record)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	ipVersionPtr := flag.Int("ip", 4, "ip version 4 or 6")
	randPtr := flag.Bool("rand", false, "sequence or random")
	uniquePtr := flag.Bool("unique", false, "unqie or identical data section")
	numPtr := flag.Int("num", 10, "number of entries")
	flag.Parse()
	writer, err := mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType: "My-ASN-DB",
			RecordSize:   24,
			IPVersion: *ipVersionPtr,
			IncludeReservedNetworks: true,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	if *ipVersionPtr == 6 {
		insert_ipv6_rand(*numPtr, writer)
	} else {
		if *randPtr {
			insert_ipv4_rand(*numPtr, writer)
		} else {
			insert_ipv4_in_sequence(*numPtr, *uniquePtr, writer)
		}
	}
	// insert_ipv6_rand(10, writer)
	// insert_ipv4_rand(10, writer)
	// insert_ipv4_in_sequence(10, true, writer)
	// insert_ipv4_in_sequence(10, false, writer)

	fh, err := os.Create("out.mmdb")
	if err != nil {
		log.Fatal(err)
	}

	_, err = writer.WriteTo(fh)
	if err != nil {
		log.Fatal(err)
	}
}
