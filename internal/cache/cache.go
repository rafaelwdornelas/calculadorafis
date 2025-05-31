package cache

import (
	"strings"
	"sync"
	"time"
)

// Instância global do cache
var (
	instance     *Cache
	instanceOnce sync.Once
)

// Item representa um item no cache com timestamp de expiração
type Item struct {
	Value      interface{}
	Expiration int64
}

// Expirado verifica se o item está expirado
func (item Item) Expirado() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

// Cache implementa um cache em memória com expiração
type Cache struct {
	items             map[string]Item
	mu                sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	stopCleanup       chan bool
}

// GetInstance retorna a instância única do cache
func GetInstance(defaultExpiration, cleanupInterval time.Duration) *Cache {
	instanceOnce.Do(func() {
		instance = newCache(defaultExpiration, cleanupInterval)
	})
	return instance
}

// newCache cria um novo cache com tempo de expiração padrão e intervalo de limpeza
func newCache(defaultExpiration, cleanupInterval time.Duration) *Cache {
	cache := &Cache{
		items:             make(map[string]Item),
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
		stopCleanup:       make(chan bool),
	}

	// Iniciar o processo de limpeza periódica se houver um intervalo definido
	if cleanupInterval > 0 {
		go cache.iniciarLimpeza()
	}

	return cache
}

// Set adiciona um item ao cache
func (c *Cache) Set(key string, value interface{}, d time.Duration) {
	var expiration int64

	if d == 0 {
		d = c.defaultExpiration
	}

	if d > 0 {
		expiration = time.Now().Add(d).UnixNano()
	}

	c.mu.Lock()
	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
	}
	c.mu.Unlock()
}

// Get obtém um item do cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}

	// Verificar se o item expirou
	if item.Expirado() {
		c.mu.RUnlock()
		return nil, false
	}

	c.mu.RUnlock()
	return item.Value, true
}

// Delete remove um item do cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

// Limpar remove todos os itens do cache
func (c *Cache) Limpar() {
	c.mu.Lock()
	c.items = make(map[string]Item)
	c.mu.Unlock()
}

// LimparItensExpirados remove todos os itens expirados do cache
func (c *Cache) LimparItensExpirados() {
	now := time.Now().UnixNano()
	c.mu.Lock()
	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			delete(c.items, k)
		}
	}
	c.mu.Unlock()
}

// iniciarLimpeza inicia a rotina de limpeza de itens expirados
func (c *Cache) iniciarLimpeza() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.LimparItensExpirados()
		case <-c.stopCleanup:
			return
		}
	}
}

// Parar interrompe a rotina de limpeza
func (c *Cache) Parar() {
	c.stopCleanup <- true
}

// StatusCache retorna estatísticas sobre o cache
func (c *Cache) StatusCache() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	agora := time.Now().UnixNano()
	totalItens := len(c.items)
	itensExpirados := 0
	itensPorPrefixo := make(map[string]int)

	// Inicializar o contador para prefixo quote_
	itensPorPrefixo["quote_"] = 0

	for k, v := range c.items {
		if v.Expiration > 0 && agora > v.Expiration {
			itensExpirados++
		}

		// Contar itens por prefixo
		if strings.HasPrefix(k, "quote_") {
			itensPorPrefixo["quote_"]++
		}
	}

	return map[string]interface{}{
		"total_itens":       totalItens,
		"itens_expirados":   itensExpirados,
		"itens_ativos":      totalItens - itensExpirados,
		"itens_por_prefixo": itensPorPrefixo,
	}
}

// ListarChavesComPrefixo retorna todas as chaves que começam com o prefixo especificado
func (c *Cache) ListarChavesComPrefixo(prefixo string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var keys []string
	for k := range c.items {
		if strings.HasPrefix(k, prefixo) {
			// Remover o prefixo para obter apenas o ticker
			ticker := strings.TrimPrefix(k, prefixo)
			keys = append(keys, ticker)
		}
	}

	return keys
}
