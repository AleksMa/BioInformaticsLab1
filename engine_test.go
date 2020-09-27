package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type EngineTestSuite struct {
	suite.Suite
	nw *NeedlemanWunsch
}

func (s *EngineTestSuite) TestDNA() {
	for _, test := range []struct {
		seq1    string
		seq2    string
		result1 string
		result2 string
		score   int
	}{
		{
			"AATCG",
			"AACG",
			"AATCG",
			"AA-CG",
			10,
		},
		{
			"GTACAACGTTA",
			"AATCGTAGCGA",
			"GTACAACGTTA",
			"AATCGTAGCGA",
			-17,
		},
		{
			"GCGCGTGCGCGGAAGGAGCCAAGGTGAAGTTGTAGCAGTGTGTCAGAAGAGGTGCGTGGCACCATGCTGTCCCCCGAGGCGGAGCGGGTGCTGCGGTACCTGGTCGAAGTAGAGGAGTTG",
			"GACTTGTGGAACCTACTTCCTGAAAATAACCTTCTGTCCTCCGAGCTCTCCGCACCCGTGGATGACCTGCTCCCGTACACAGATGTTGCCACCTGGCTGGATGAATGTCCGAATGAAGCG",
			"GCGCGTGCGCGGAAGGAGCCAAGGTGAAGTTGTAGCAGTGTGTCAGAAGAGGTGCGTGGCA-CCAT-GCTGTCCCCCGAGGCGGA-GCGGGTGCTG-C-GGTACCTGGTCGAA-GT-AG-AGGAGTTG",
			"G-AC-T-TGTGGAA-CCTACTTCCTGAA--AATAACCTTCTGTCCTCCGAGCT-CTCCGCACCCGTGGATGACCTGC-TCCCGTACACAGATGTTGCCACCTGGCTGGATGAATGTCCGAATGAAGCG",
			-41,
		},
	} {
		s.nw = NewNeedlemanWunsch(&Sequence{Value: test.seq1}, &Sequence{Value: test.seq2}, DNAFull, -10)
		r1, r2, sc := s.nw.Solve()
		s.Equal(r1, test.result1)
		s.Equal(r2, test.result2)
		s.Equal(sc, test.score)
	}
}

func (s *EngineTestSuite) TestProteins() {
	for _, test := range []struct {
		seq1    string
		seq2    string
		result1 string
		result2 string
		score   int
	}{
		{
			"SPETVIHSGWVIWRELFSHWPDQCKLLFGDWFAWIHWTYLVYYSAGPPCQGQSDIVVMMQKKLRTNFCQCYKYWYQ",
			"SPSDQFFTVIHSCLYWVIWRDLMSHLFMNGAAIDIHWTWDSIAIGPPLVYPIEEVFAGPSTIVVMMQKMLRTNFCQCYKPWYQ",
			"SP--E--TVIHS--GWVIWRELFSH-WPDQCKL-LFGDWFAWIHWTYLVYYSAGPPCQGQSDIVVMMQKKLRTNFCQCYKYWYQ",
			"SPSDQFFTVIHSCLYWVIWRDLMSHLFMNGAAIDIHWTWDSIAIGPPLV-YPIEEVFAGPSTIVVMMQKMLRTNFCQCYKPWYQ",
			116,
		},
		{
			"FAWIHWTYLVYYSAGPPCQGQSDNFCQCYKYWYQQC",
			"KMLRTNFCQCYDSIAIGPPLVYPIEEVFAGPSTIVVMMQ",
			"FAWIHWTYL-VYYS-A-GPPCQGQSDN-FCQCYKYWYQQC",
			"-KMLRTNFCQCYDSIAIGPPLVYPIEEVFAGPSTIVVMMQ",
			-36,
		},
		{
			"ARDCDWVKMF",
			"IGKWVKDN",
			"ARDCDWVKMF",
			"I-G-KWVKDN",
			-9,
		},
	} {
		s.nw = NewNeedlemanWunsch(&Sequence{Value: test.seq1}, &Sequence{Value: test.seq2}, Blosum62, -10)
		r1, r2, sc := s.nw.Solve()
		s.Equal(r1, test.result1)
		s.Equal(r2, test.result2)
		s.Equal(sc, test.score)
	}
}

func TestVkAudienceCountersSuite(t *testing.T) {
	suite.Run(t, new(EngineTestSuite))
}
