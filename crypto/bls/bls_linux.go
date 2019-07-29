// +build linux darwin

package bls

/*
#cgo CFLAGS:-DMCLBN_FP_UNIT_SIZE=6
#cgo LDFLAGS:-lbls384 -lgmpxx -lstdc++ -lgmp -lcrypto
#include <bls/bls.h>
*/
import "C"
import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unsafe"

	"github.com/PlatONnetwork/PlatON-Go/crypto"
)

// Init --
// call this function before calling all the other operations
// this function is not thread safe
func Init(curve int) error {
	err := C.blsInit(C.int(curve), C.MCLBN_FP_UNIT_SIZE)
	if err != 0 {
		return fmt.Errorf("ERR Init curve=%d", curve)
	}
	return nil
}

// ID --
type ID struct {
	v Fr
}

// getPointer --
func (id *ID) getPointer() (p *C.blsId) {
	// #nosec
	return (*C.blsId)(unsafe.Pointer(id))
}

// GetLittleEndian --
func (id *ID) GetLittleEndian() []byte {
	return id.v.Serialize()
}

// SetLittleEndian --
func (id *ID) SetLittleEndian(buf []byte) error {
	return id.v.SetLittleEndian(buf)
}

// GetHexString --
func (id *ID) GetHexString() string {
	return id.v.GetString(16)
}

// GetDecString --
func (id *ID) GetDecString() string {
	return id.v.GetString(10)
}

// SetHexString --
func (id *ID) SetHexString(s string) error {
	return id.v.SetString(s, 16)
}

// SetDecString --
func (id *ID) SetDecString(s string) error {
	return id.v.SetString(s, 10)
}

// IsEqual --
func (id *ID) IsEqual(rhs *ID) bool {
	return id.v.IsEqual(&rhs.v)
}

// SecretKey --
type SecretKey struct {
	v Fr
}

func LoadBLS(file string) (*SecretKey, error) {
	buf := make([]byte, 64)
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	if _, err := io.ReadFull(fd, buf); err != nil {
		return nil, err
	}
	var sec SecretKey
	key, err := hex.DecodeString(string(buf))
	if err != nil {
		return nil, err
	}
	err = sec.SetLittleEndian(key)
	return &sec, err
}

// getPointer --
func (sec *SecretKey) getPointer() (p *C.blsSecretKey) {
	// #nosec
	return (*C.blsSecretKey)(unsafe.Pointer(sec))
}

// GetLittleEndian --
func (sec *SecretKey) GetLittleEndian() []byte {
	return sec.v.Serialize()
}

// SetLittleEndian --
func (sec *SecretKey) SetLittleEndian(buf []byte) error {
	return sec.v.SetLittleEndian(buf)
}

// GetHexString --
func (sec *SecretKey) GetHexString() string {
	return sec.v.GetString(16)
}

// GetDecString --
func (sec *SecretKey) GetDecString() string {
	return sec.v.GetString(10)
}

// SetHexString --
func (sec *SecretKey) SetHexString(s string) error {
	return sec.v.SetString(s, 16)
}

// SetDecString --
func (sec *SecretKey) SetDecString(s string) error {
	return sec.v.SetString(s, 10)
}

// IsEqual --
func (sec *SecretKey) IsEqual(rhs *SecretKey) bool {
	return sec.v.IsEqual(&rhs.v)
}

// SetByCSPRNG --
func (sec *SecretKey) SetByCSPRNG() {
	sec.v.SetByCSPRNG()
}

// Add --
func (sec *SecretKey) Add(rhs *SecretKey) {
	FrAdd(&sec.v, &sec.v, &rhs.v)
}

// Mul --
func (sec *SecretKey) Mul(rhs *SecretKey) {
	FrMul(&sec.v, &sec.v, &rhs.v)
}

// Sub --
func (sec *SecretKey) Sub(rhs *SecretKey) {
	FrSub(&sec.v, &sec.v, &rhs.v)
}

// GetMasterSecretKey --
func (sec *SecretKey) GetMasterSecretKey(k int) (msk []SecretKey) {
	msk = make([]SecretKey, k)
	msk[0] = *sec
	for i := 1; i < k; i++ {
		msk[i].SetByCSPRNG()
	}
	return msk
}

// GetMasterPublicKey --
func GetMasterPublicKey(msk []SecretKey) (mpk []PublicKey) {
	n := len(msk)
	mpk = make([]PublicKey, n)
	for i := 0; i < n; i++ {
		mpk[i] = *msk[i].GetPublicKey()
	}
	return mpk
}

// Set --
func (sec *SecretKey) Set(msk []SecretKey, id *ID) error {
	// #nosec
	return FrEvaluatePolynomial(&sec.v, *(*[]Fr)(unsafe.Pointer(&msk)), &id.v)
}

// Recover --
func (sec *SecretKey) Recover(secVec []SecretKey, idVec []ID) error {
	// #nosec
	return FrLagrangeInterpolation(&sec.v, *(*[]Fr)(unsafe.Pointer(&idVec)), *(*[]Fr)(unsafe.Pointer(&secVec)))
}

// GetPop --
func (sec *SecretKey) GetPop() (sign *Sign) {
	sign = new(Sign)
	C.blsGetPop(sign.getPointer(), sec.getPointer())
	return sign
}

// PublicKey --
type PublicKey struct {
	v G2
}

func (pub *PublicKey) getQ() (p *C.blsPublicKey) {
	// #nosec
	return (*C.blsPublicKey)(unsafe.Pointer(pub))
}

// getPointer --
func (pub *PublicKey) getPointer() (p *C.blsPublicKey) {
	// #nosec
	return (*C.blsPublicKey)(unsafe.Pointer(pub))
}

// Serialize --
func (pub *PublicKey) Serialize() []byte {
	return pub.v.Serialize()
}

// Deserialize --
func (pub *PublicKey) Deserialize(buf []byte) error {
	return pub.v.Deserialize(buf)
}

// GetHexString --
func (pub *PublicKey) GetHexString() string {
	return pub.v.GetString(16)
}

// SetHexString --
func (pub *PublicKey) SetHexString(s string) error {
	return pub.v.SetString(s, 16)
}

// IsEqual --
func (pub *PublicKey) IsEqual(rhs *PublicKey) bool {
	return pub.v.IsEqual(&rhs.v)
}

func (pub *PublicKey) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%x", pub.Serialize())), nil
}

func (pub *PublicKey) UnmarshalText(text []byte) error {
	key, err := hex.DecodeString(string(text))
	if err != nil {
		return err
	}
	return pub.Deserialize(key)
}

// Add --
func (pub *PublicKey) Add(rhs *PublicKey) {
	G2Add(&pub.v, &pub.v, &rhs.v)
}

// Mul --
func (pub *PublicKey) Mul(rhs *SecretKey) {
	G2Mul(&pub.v, &pub.v, &rhs.v)
}

// Set --
func (pub *PublicKey) Set(mpk []PublicKey, id *ID) error {
	// #nosec
	return G2EvaluatePolynomial(&pub.v, *(*[]G2)(unsafe.Pointer(&mpk)), &id.v)
}

// Recover --
func (pub *PublicKey) Recover(pubVec []PublicKey, idVec []ID) error {
	// #nosec
	return G2LagrangeInterpolation(&pub.v, *(*[]Fr)(unsafe.Pointer(&idVec)), *(*[]G2)(unsafe.Pointer(&pubVec)))
}

// Sign  --
type Sign struct {
	v G1
}

// getPointer --
func (sign *Sign) getPointer() (p *C.blsSignature) {
	// #nosec
	return (*C.blsSignature)(unsafe.Pointer(sign))
}

// Serialize --
func (sign *Sign) Serialize() []byte {
	return sign.v.Serialize()
}

// Deserialize --
func (sign *Sign) Deserialize(buf []byte) error {
	return sign.v.Deserialize(buf)
}

// GetHexString --
func (sign *Sign) GetHexString() string {
	return sign.v.GetString(16)
}

// SetHexString --
func (sign *Sign) SetHexString(s string) error {
	return sign.v.SetString(s, 16)
}

// IsEqual --
func (sign *Sign) IsEqual(rhs *Sign) bool {
	return sign.v.IsEqual(&rhs.v)
}

// GetPublicKey --
func (sec *SecretKey) GetPublicKey() (pub *PublicKey) {
	pub = new(PublicKey)
	C.blsGetPublicKey(pub.getPointer(), sec.getPointer())
	return pub
}

// Sign -- Constant Time version
func (sec *SecretKey) Sign(m string) (sign *Sign) {
	sign = new(Sign)
	buf := []byte(m)
	// #nosec
	C.blsSign(sign.getPointer(), sec.getPointer(), unsafe.Pointer(&buf[0]), C.size_t(len(buf)))
	return sign
}

// Add --
func (sign *Sign) Add(rhs *Sign) {
	C.blsSignatureAdd(sign.getPointer(), rhs.getPointer())
}

// Recover --
func (sign *Sign) Recover(signVec []Sign, idVec []ID) error {
	// #nosec
	return G1LagrangeInterpolation(&sign.v, *(*[]Fr)(unsafe.Pointer(&idVec)), *(*[]G1)(unsafe.Pointer(&signVec)))
}

// Verify --
func (sign *Sign) Verify(pub *PublicKey, m string) bool {
	buf := []byte(m)
	// #nosec
	return C.blsVerify(sign.getPointer(), pub.getPointer(), unsafe.Pointer(&buf[0]), C.size_t(len(buf))) == 1
}

// VerifyPop --
func (sign *Sign) VerifyPop(pub *PublicKey) bool {
	return C.blsVerifyPop(sign.getPointer(), pub.getPointer()) == 1
}

// DHKeyExchange --
func DHKeyExchange(sec *SecretKey, pub *PublicKey) (out PublicKey) {
	C.blsDHKeyExchange(out.getPointer(), sec.getPointer(), pub.getPointer())
	return out
}

//add@20190716
//get G2
func GetGeneratorOfG2() (pub *PublicKey) {
	pub = new(PublicKey)
	C.blsGetGeneratorOfG2(pub.getPointer())
	return pub
}

// PubBatchAdd --
func PubkeyBatchAdd(pkVec []PublicKey) (pub PublicKey) {
	var pk PublicKey
	for i := 0; i < len(pkVec); i++ {
		pk.Add(&pkVec[i])
	}
	fmt.Printf("pk=%s\n", pk.GetHexString())
	return pk
}

// SecBatchAdd --
func SeckeyBatchAdd(secVec []SecretKey) (sec SecretKey) {
	var sk SecretKey
	for i := 0; i < len(secVec); i++ {
		sk.Add(&secVec[i])
	}
	fmt.Printf("sec=%s\n", sk.GetHexString())
	return sk
}

func AggregateSign(sigVec []Sign) (sig Sign) {
	var sign Sign
	for i := 0; i < len(sigVec); i++ {
		sign.Add(&sigVec[i])
	}
	fmt.Printf("sig=%s\n", sign.GetHexString())
	return sign
}

func GTBatchMul(eVec []GT) (e GT) {
	var e1, e2 GT
	e1 = eVec[0]
	for j := 1; j < len(eVec); j++ {
		e2 = eVec[j]
		GTMul(&e2, &e1, &e2)
		e1 = e2
	}
	fmt.Printf("e=%s\n", e2.GetString(16))
	return e2
}

func GTBatchAdd(eVec []GT) (e GT) {
	var e1, e2 GT
	e1 = eVec[0]
	for j := 1; j < len(eVec); j++ {
		e2 = eVec[j]
		GTAdd(&e2, &e1, &e2)
		e1 = e2
	}
	fmt.Printf("e=%s\n", e2.GetString(16))
	return e2
}

func MsgsToHashToG1(mVec []string) ([]Sign, error) {
	n := len(mVec)
	p_Hm := make([]Sign, n)
	for i := 0; i < n; i++ {
		err := p_Hm[i].v.HashAndMapTo([]byte(mVec[i]))
		if err != nil {
			return []Sign{}, err
		}
		fmt.Printf("p_Hm=%s\n", p_Hm[i].GetHexString())
	}
	return p_Hm, nil
}

func BatchVerifySameMsg(curve int, msg string, pkVec []PublicKey, sign Sign) error {
	err := Init(curve)
	if err != nil {
		return err
	}
	/*	if len(pkVec) != len(signVec) {
		return errors.New("sig/pub length not equal")
	}*/
	var pk PublicKey
	//	var sig Sign
	for i := 0; i < len(pkVec); i++ {
		pk.Add(&pkVec[i])
		//		sig.Add(&signVec[i])
	}
	fmt.Printf("pk=%s\n", pk.GetHexString())
	//	fmt.Printf("sig=%s\n", sig.GetHexString())
	if !sign.Verify(&pk, msg) {
		return errors.New("signature verification failed")
	}
	return nil
}

func BatchVerifyDistinctMsg(curve int, pkVec []PublicKey, msgVec []Sign, sign Sign) error {
	err := Init(curve)
	if err != nil {
		return err
	}
	if len(pkVec) != len(msgVec) {
		return errors.New("pub/msg length not equal")
	}
	/*var sig Sign
	for i := 0; i < len(pkVec); i++ {
		sig.Add(&signVec[i])
	}*/
	//	fmt.Printf("sig=%s\n", sig.GetHexString())
	P := GetGeneratorOfG2()
	fmt.Printf("G=%s\n", P.GetHexString())
	var e, e1, e2 GT
	Pairing(&e, &(sign.v), &(P.v))
	fmt.Printf("e=%s\n", e.GetString(16))

	n := len(msgVec)
	Pairing(&e1, &(msgVec[0].v), &(pkVec[0].v))
	for j := 1; j < n; j++ {
		Pairing(&e2, &(msgVec[j].v), &(pkVec[j].v))
		GTMul(&e2, &e1, &e2)
		e1 = e2
	}
	fmt.Printf("e2=%s\n", e2.GetString(16))
	if !e.IsEqual(&e2) {
		errors.New("not equal pairing\n")
	}
	return nil
}

func SynthSameMsg(curve int, pkVec []PublicKey, mVec []string) ([]PublicKey, []string, []string, error) {
	err := Init(curve)
	if err != nil {
		return nil, nil, nil, err
	}
	pubMap := make(map[string]PublicKey)
	mark := make(map[string]string)
	for i := 0; i < len(mVec); i++ {
		pub, ok := pubMap[mVec[i]]
		if ok {
			pub.Add(&pkVec[i])
			pubMap[mVec[i]] = pub
			mark[mVec[i]] = mark[mVec[i]] + fmt.Sprintf(",%d", i)
		} else {
			pubMap[mVec[i]] = pkVec[i]
			mark[mVec[i]] = fmt.Sprintf("%d", i)
		}
	}
	n := len(pubMap)
	newPkVec := make([]PublicKey, n)
	newMVec := make([]string, n)
	var j int = 0
	for k, v := range pubMap {
		newMVec[j] = k
		newPkVec[j] = v
		j++
	}
	var index []string
	for _, result := range mark {
		if strings.Contains(result, ",") {
			index = append(index, result)
		}
	}
	return newPkVec, newMVec, index, nil
}

// IsValid --
func G2IsValid(rhs *PublicKey) bool {
	return C.mclBnG2_isValid((&rhs.v).getPointer()) == 1
}

func Schnorr_test(curve int, r, c SecretKey, G, V, P PublicKey) error {
	err := Init(curve)
	if err != nil {
		return err
	}
	fmt.Printf("r=%s\n", r.GetHexString())
	fmt.Printf("c=%s\n", c.GetHexString())
	fmt.Printf("G=%s\n", G.GetHexString())
	fmt.Printf("V=%s\n", V.GetHexString())
	fmt.Printf("P=%s\n", P.GetHexString())
	if !G2IsValid(&P) {
		return errors.New("P isnot valid")
	}
	Pr := G
	Pr.Mul(&r)
	Pc := P
	Pc.Mul(&c)
	Psum := Pr
	Psum.Add(&Pc)
	fmt.Printf("Psum=%s\n", Psum.GetHexString())
	if !V.IsEqual(&Psum) {
		return errors.New("V = G*[r] + P*[c] not equal")
	}
	return nil
}

type Proof struct {
	C, R SecretKey
}

func SchnorrNIZKProve(curve int, sec SecretKey) (*Proof, error) {
	err := Init(curve)
	if err != nil {
		return nil, err
	}
	P := sec.GetPublicKey()
	fmt.Printf("P=%s\n", P.GetHexString())
	var v SecretKey
	v.SetByCSPRNG()
	V := v.GetPublicKey()
	fmt.Printf("V=%s\n", V.GetHexString())
	G := GetGeneratorOfG2()
	fmt.Printf("G=%s\n", G.GetHexString())

	input1 := G.Serialize()
	input2 := P.Serialize()
	input3 := V.Serialize()
	fmt.Printf("input1=%x\n", input1)
	fmt.Printf("input2=%x\n", input2)
	fmt.Printf("input3=%x\n", input3)

	var buffer bytes.Buffer
	buffer.Write(input1)
	buffer.Write(input2)
	buffer.Write(input3)
	output := buffer.Bytes()
	h := crypto.Keccak256(output)
	fmt.Printf("h=%x\n", h)
	var c SecretKey
	err = c.SetLittleEndian(h)
	if err != nil {
		fmt.Println("Deserialize fail")
		return nil, err
	}
	temp := sec
	temp.Mul(&c)
	r := v
	r.Sub(&temp)
	fmt.Printf("r=%s\n", r.GetHexString())
	fmt.Printf("c=%s\n", c.GetHexString())
	proof := new(Proof)
	proof.C = c
	proof.R = r
	return proof, nil
}

func SchnorrNIZKVerify(curve int, proof Proof, P PublicKey) error {
	err := Init(curve)
	if err != nil {
		return err
	}
	if !G2IsValid(&P) {
		return errors.New("P isnot valid")
	}
	c := proof.C
	r := proof.R
	G := GetGeneratorOfG2()
	fmt.Printf("G=%s\n", G.GetHexString())
	//	V1 = G * r + A * c     c = H(G || pk || V’)
	var Pr PublicKey
	Pr = *G
	Pr.Mul(&r)
	Pc := P
	Pc.Mul(&c)
	V1 := Pr
	V1.Add(&Pc)
	fmt.Printf("V1=%s\n", V1.GetHexString())
	fmt.Printf("P=%s\n", P.GetHexString())

	input1 := G.Serialize()
	input2 := P.Serialize()
	input3 := V1.Serialize()
	fmt.Printf("input1=%x\n", input1)
	fmt.Printf("input2=%x\n", input2)
	fmt.Printf("input3=%x\n", input3)
	var buffer bytes.Buffer
	buffer.Write(input1)
	buffer.Write(input2)
	buffer.Write(input3)
	output := buffer.Bytes()
	h := crypto.Keccak256(output)
	fmt.Printf("h=%x\n", h)
	var c1 SecretKey
	err = c1.SetLittleEndian(h)
	if err != nil {
		return err
	}
	fmt.Printf("c1=%s\n", c1.GetHexString())
	if !c.IsEqual(&c1) {
		return errors.New("not same c = H(G || pk || V’)")
	}
	return nil
}