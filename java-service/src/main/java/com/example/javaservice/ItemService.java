package com.example.javaservice;

import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicLong;

@Service
public class ItemService {
    private final Map<Long, Item> items = new ConcurrentHashMap<>();
    private final AtomicLong idGenerator = new AtomicLong(0);

    public List<Item> list() {
        return new ArrayList<>(items.values());
    }

    public Item create(ItemRequest request) {
        Long id = idGenerator.incrementAndGet();
        Item item = new Item(id, request.getName(), request.getDescription());
        items.put(id, item);
        return item;
    }

    public Item get(Long id) {
        return items.get(id);
    }

    public Item update(Long id, ItemRequest request) {
        Item existing = items.get(id);
        if (existing == null) {
            return null;
        }

        existing.setName(request.getName());
        existing.setDescription(request.getDescription());
        return existing;
    }

    public boolean delete(Long id) {
        return items.remove(id) != null;
    }
}
