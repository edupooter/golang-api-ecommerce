Instruções do Agente — Diretrizes para implementar funcionalidades core de e‑commerce

Objetivo
- Fornecer regras claras que o agente deve seguir ao implementar carrinho de compras, cadastro de comprador, checkout e baixa de estoque.
- Observar DDD, SOLID, DRY, KISS e YAGNI; aplicar design patterns quando justificável.

Aplicabilidade
- Valem para o código em internal/, cmd/, e testes. Prefira organizar domínio, serviços, repositórios e handlers conforme indicado abaixo.

Princípios (o que aplicar na prática)
- DDD: modele entidades e agregados (Product, Customer, Cart, CartItem, Order). Coloque lógica de negócio no domínio/serviços, não nos handlers.
- SOLID: single responsibility por camada; dependências por interfaces; injeção via construtores; não quebre LSP ao alterar contratos.
- DRY: extrair helpers para leitura/escrita JSON, mapeamento de erros HTTP e validações reutilizáveis.
- KISS: prefira soluções simples e diretas; evite abstrações desnecessárias.
- YAGNI: não implemente patterns complexos (event bus, CQRS, eventual consistency) sem necessidade clara.

Estrutura recomendada (hexagonal/layers)
- Domain (modelos e regras): internal/domain ou internal/model — entidades com métodos que preservam invariantes.
- Ports (interfaces): internal/repo — defina interfaces de repositório (ProductRepository, OrderRepository, CustomerRepository).
- Adapters (implementações): internal/repo/sqlite.go e internal/repo/memory.go — implementações dos ports.
- Application / Services: internal/service — ProductService, CartService, OrderService; aqui ficam validações e orquestração.
- Delivery: internal/handler — handlers HTTP leves (parsing, headers, status) que delegam para services.
- Infra: cmd/server e internal/server (router, utilitários HTTP).

Modelagem do domínio (diretrizes)
- Entidades encapsulam comportamento: ex.: Cart.AddItem(item), Product.CanReserve(qty), Order.CalculateTotal().
- Value objects quando necessário (Money); se não houver multi-currency, considerar YAGNI.
- Agregados: Order é agregado que contém items e total; Cart pode ser agregado temporário até checkout.

Fluxo de checkout (resumo técnico)
1. Validar carrinho e cliente no service.
2. Iniciar transação (DB) quando possível.
3. Para cada item: decrementar estoque de forma atômica:
   - SQL robusto: UPDATE products SET stock = stock - ? WHERE id = ? AND stock >= ?; verificar RowsAffected.
   - Em repositório in-memory: proteger com mutex.
4. Se qualquer decremento falhar: rollback e retornar ErrInsufficientStock (HTTP 409).
5. Criar Order e persistir dentro da transação.
6. Commit e retornar pedido criado.

Concorrência e consistência
- Em DB relacional use atualização atômica (consulta/UPDATE condicional) ou locking apropriado.
- Como alternativa use coluna de versão (optimistic locking) se necessário.
- Para SQLite, usar transação e UPDATE condicional; para in-memory, use mutex.

Mapeamento de erros para HTTP
- ErrNotFound -> 404
- ErrInvalid -> 400
- ErrInsufficientStock -> 409
- Erro interno não esperado -> 500
- Centralize esse mapeamento em internal/server/errors.go e use helpers de resposta.

Testing
- Unit tests para domain (regras) e services (negócio) com mocks para repositórios.
- Testes de integração usando sqlite em arquivo temporário ou repo em memória.

Patterns recomendados (quando justificados)
- Repository, Unit of Work (para transações), Factory/Builder para criação de Orders complexos, Strategy para regras de preço/promos.
- Evite Event Sourcing/CQRS até que a necessidade fique clara (YAGNI).

Boas práticas de implementação
- Handlers: apenas parse, validação sintática mínima e delegação.
- Services: validação semântica e orquestração transacional.
- Repositories: contratos simples, retornando erros padronizados.
- Helpers: internal/server/response.go com WriteJSON/ReadJSON; internal/server/errors.go com MapErrorToStatus.


Exemplos de assinatura de serviço
- func (s *OrderService) Checkout(ctx context.Context, customerID int64, cart *domain.Cart) (*domain.Order, error)
- func (s *CartService) AddItem(ctx context.Context, cartID string, item domain.CartItem) error

Decisões do projeto (respondidas)
- Arquitetura: migrar para Hexagonal Architecture (ports/adapters/services/handlers).
- Concorrência: controle via UPDATE condicional (atômico) para bancos; repositório em memória protegido por mutex.
- Transações: iniciar com transações simples (Unit-of-Work onde aplicável), evoluir conforme necessidade.
- Scaffolding: scaffold inicial (domain, services, handlers, testes) foi gerado e integrado.

Próximos passos sugeridos
- Implementar `OrderRepository` e `CustomerRepository` (memória + SQLite) e registrar rota `/checkout`.
- Extrair helpers HTTP (`internal/server/response.go`) e mapeamento de erros (`internal/server/errors.go`).
- Adicionar testes de integração para fluxo de checkout com SQLite.

Notas finais
- Aplique YAGNI: comece simples e só introduza patterns avançados quando os requisitos exigirem.
- Mantenha contratos (interfaces) pequenos e estáveis para facilitar refatorações e testes.

Fim das instruções.
