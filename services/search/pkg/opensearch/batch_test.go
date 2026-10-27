package opensearch_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/opencloud-eu/opencloud/services/search/internal/opensearchtest"
	"github.com/opencloud-eu/opencloud/services/search/pkg/opensearch"
)

var _ = Describe("Batch", func() {
	Describe("Push", func() {
		var tc *opensearchtest.TestClient

		BeforeEach(func() {
			tc = opensearchtest.NewDefaultTestClient(GinkgoTB(), defaultConfig.Engine.OpenSearch.Client)
		})

		It("reports the documents the bulk API rejected", func() {
			const indexName = "opencloud-test-batch-push-rejected"
			tc.Require.IndicesReset([]string{indexName})
			deleteIndexOnCleanup(tc, indexName)

			// Name is a string, mapping it as a long makes every document fail to parse.
			tc.Require.IndicesCreate(indexName, strings.NewReader(`{"mappings":{"properties":{"Name":{"type":"long"}}}}`))

			batch, err := opensearch.NewBatch(tc.Client(), indexName, 10)
			Expect(err).ToNot(HaveOccurred())

			document := opensearchtest.Testdata.Resources.File
			Expect(batch.Upsert(document.ID, document)).To(Succeed())

			err = batch.Push()
			Expect(err).To(MatchError(ContainSubstring(document.ID)))
			Expect(err).To(MatchError(ContainSubstring("mapper_parsing_exception")))
			tc.Require.IndicesCount([]string{indexName}, nil, 0)
		})

		It("pushes the documents the bulk API accepted", func() {
			const indexName = "opencloud-test-batch-push-accepted"
			tc.Require.IndicesReset([]string{indexName})
			deleteIndexOnCleanup(tc, indexName)

			tc.Require.IndicesCreate(indexName, strings.NewReader(opensearch.IndexManagerLatest.String()))

			batch, err := opensearch.NewBatch(tc.Client(), indexName, 10)
			Expect(err).ToNot(HaveOccurred())

			document := opensearchtest.Testdata.Resources.File
			Expect(batch.Upsert(document.ID, document)).To(Succeed())
			Expect(batch.Push()).To(Succeed())

			tc.Require.IndicesCount([]string{indexName}, nil, 1)
		})
	})
})
