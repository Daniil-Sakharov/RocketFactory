//go:build integration

package integration

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/inventory/v1"
)

var _ = Describe("InventoryService", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventoryV1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Даем Docker время настроить port mapping
		// Приложение готово внутри контейнера, но mapped port может быть еще не доступен с хоста
		// В CI/CD среде может потребоваться больше времени
		time.Sleep(3 * time.Second)

		// Создаём gRPC клиент
		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventoryV1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		// Чистим коллекцию после теста
		err := env.ClearPartsCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку коллекции parts")

		cancel()
	})

	Describe("ListParts", func() {
		BeforeEach(func() {
			// Вставляем тестовые детали в коллекцию
			err := env.InsertTestParts(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовых деталей в MongoDB")
		})

		It("должен успешно возвращать список всех деталей без фильтра", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).ToNot(BeEmpty(), "список деталей не должен быть пустым")
			Expect(len(resp.GetParts())).To(BeNumerically(">", 0))

			// Проверяем структуру первой детали
			if len(resp.GetParts()) > 0 {
				part := resp.GetParts()[0]
				Expect(part.Uuid).ToNot(BeEmpty(), "UUID детали не должен быть пустым")
				Expect(part.Name).ToNot(BeEmpty(), "Name детали не должно быть пустым")
				Expect(part.Category).ToNot(Equal(inventoryV1.Category_CATEGORY_UNSPECIFIED))
			}
		})

		It("должен успешно фильтровать детали по UUID", func() {
			// Сначала получаем список всех деталей
			allPartsResp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{})
			Expect(err).ToNot(HaveOccurred())
			Expect(allPartsResp.GetParts()).ToNot(BeEmpty())

			// Берем UUID первой детали
			firstPartUUID := allPartsResp.GetParts()[0].Uuid

			// Фильтруем по этому UUID
			filteredResp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Uuids: []string{firstPartUUID},
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(filteredResp.GetParts()).To(HaveLen(1))
			Expect(filteredResp.GetParts()[0].Uuid).To(Equal(firstPartUUID))
		})

		It("должен успешно фильтровать детали по категории", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Categories: []inventoryV1.Category{inventoryV1.Category_CATEGORY_ENGINE},
				},
			})

			Expect(err).ToNot(HaveOccurred())

			// Проверяем, что все возвращенные детали имеют категорию ENGINE
			for _, part := range resp.GetParts() {
				Expect(part.Category).To(Equal(inventoryV1.Category_CATEGORY_ENGINE))
			}
		})
	})

	Describe("GetPart", func() {
		var testPartUUID string

		BeforeEach(func() {
			// Вставляем одну тестовую деталь
			var err error
			testPartUUID, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовой детали в MongoDB")
		})

		It("должен успешно возвращать деталь по UUID", func() {
			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: testPartUUID,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetPart()).ToNot(BeNil())
			Expect(resp.GetPart().Uuid).To(Equal(testPartUUID))
			Expect(resp.GetPart().Name).ToNot(BeEmpty())
			Expect(resp.GetPart().Description).ToNot(BeEmpty())
			Expect(resp.GetPart().Price).To(BeNumerically(">", 0))
			Expect(resp.GetPart().StockQuantity).To(BeNumerically(">=", 0))
			Expect(resp.GetPart().Category).ToNot(Equal(inventoryV1.Category_CATEGORY_UNSPECIFIED))
			Expect(resp.GetPart().GetDimensions()).ToNot(BeNil())
			Expect(resp.GetPart().GetManufacturer()).ToNot(BeNil())
			Expect(resp.GetPart().GetCreatedAt()).ToNot(BeNil())
		})

		It("должен вернуть ошибку при запросе несуществующей детали", func() {
			_, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: "non-existent-uuid",
			})

			Expect(err).To(HaveOccurred(), "ожидали ошибку при запросе несуществующей детали")
		})
	})

	Describe("Полный сценарий работы с инвентарем", func() {
		It("должен поддерживать полный сценарий: вставка и получение деталей", func() {
			// 1. Вставляем несколько тестовых деталей
			err := env.InsertTestParts(ctx)
			Expect(err).ToNot(HaveOccurred())

			// 2. Получаем список всех деталей
			listResp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{})
			Expect(err).ToNot(HaveOccurred())
			Expect(listResp.GetParts()).ToNot(BeEmpty())

			// 3. Берем UUID первой детали
			firstPartUUID := listResp.GetParts()[0].Uuid

			// 4. Получаем деталь по UUID
			getResp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: firstPartUUID,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(getResp.GetPart().Uuid).To(Equal(firstPartUUID))

			// 5. Фильтруем по категории детали
			partCategory := getResp.GetPart().Category
			filteredResp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Categories: []inventoryV1.Category{partCategory},
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(filteredResp.GetParts()).ToNot(BeEmpty())

			// Проверяем, что все отфильтрованные детали имеют нужную категорию
			for _, part := range filteredResp.GetParts() {
				Expect(part.Category).To(Equal(partCategory))
			}
		})
	})
})
