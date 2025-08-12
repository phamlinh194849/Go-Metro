const redis = require('redis');

// Cấu hình Redis connection
const redisConfig = {
  host: "localhost",
  port: "6379",
};

// Tạo Redis client
const client = redis.createClient({
  socket: {
    host: redisConfig.host,
    port: parseInt(redisConfig.port)
  },
  password: redisConfig.password
});

// Hàm clear toàn bộ cache
async function clearAllCache() {
  try {
    await client.connect();
    console.log('Đã kết nối Redis thành công');
    
    // Xóa toàn bộ database hiện tại
    await client.flushDb();
    console.log('✅ Đã clear toàn bộ cache thành công!');
    
  } catch (error) {
    console.error('❌ Lỗi khi clear cache:', error.message);
  } finally {
    await client.disconnect();
    console.log('Đã ngắt kết nối Redis');
  }
}

// Hàm clear cache theo pattern
async function clearCacheByPattern(pattern) {
  try {
    await client.connect();
    console.log('Đã kết nối Redis thành công');
    
    // Tìm tất cả keys theo pattern
    const keys = await client.keys(pattern);
    console.log(`Tìm thấy ${keys.length} keys với pattern: ${pattern}`);
    
    if (keys.length > 0) {
      // Xóa các keys
      await client.del(keys);
      console.log(`✅ Đã xóa ${keys.length} cache keys!`);
    } else {
      console.log('Không có keys nào để xóa');
    }
    
  } catch (error) {
    console.error('❌ Lỗi khi clear cache:', error.message);
  } finally {
    await client.disconnect();
    console.log('Đã ngắt kết nối Redis');
  }
}

// Hàm clear cache theo danh sách keys cụ thể
async function clearSpecificKeys(keys) {
  try {
    await client.connect();
    console.log('Đã kết nối Redis thành công');
    
    // Kiểm tra keys có tồn tại không
    const existingKeys = [];
    for (const key of keys) {
      const exists = await client.exists(key);
      if (exists) {
        existingKeys.push(key);
      }
    }
    
    console.log(`Tìm thấy ${existingKeys.length}/${keys.length} keys tồn tại`);
    
    if (existingKeys.length > 0) {
      await client.del(existingKeys);
      console.log(`✅ Đã xóa ${existingKeys.length} keys:`, existingKeys);
    } else {
      console.log('Không có keys nào để xóa');
    }
    
  } catch (error) {
    console.error('❌ Lỗi khi clear cache:', error.message);
  } finally {
    await client.disconnect();
    console.log('Đã ngắt kết nối Redis');
  }
}

// Hàm liệt kê tất cả keys (để debug)
async function listAllKeys() {
  try {
    await client.connect();
    console.log('Đã kết nối Redis thành công');
    
    const keys = await client.keys('*');
    console.log(`Tổng số keys: ${keys.length}`);
    
    if (keys.length > 0) {
      console.log('Danh sách keys:');
      keys.forEach((key, index) => {
        console.log(`${index + 1}. ${key}`);
      });
    }
    
  } catch (error) {
    console.error('❌ Lỗi khi list keys:', error.message);
  } finally {
    await client.disconnect();
    console.log('Đã ngắt kết nối Redis');
  }
}

// Các cách sử dụng:

// 1. Clear toàn bộ cache
// clearAllCache();

// 2. Clear cache theo pattern (ví dụ: tất cả keys bắt đầu bằng "user:")
// clearCacheByPattern('user:*');

// 3. Clear cache theo pattern khác
// clearCacheByPattern('session:*');
// clearCacheByPattern('cache:*');

// 4. Clear specific keys
// clearSpecificKeys(['user:123', 'session:abc', 'cache:data']);

// 5. List tất cả keys để xem
// listAllKeys();

// Export các functions để sử dụng ở nơi khác
module.exports = {
  clearAllCache,
  clearCacheByPattern,
  clearSpecificKeys,
  listAllKeys
};

// Nếu chạy file này trực tiếp
if (require.main === module) {
  console.log('Redis Cache Clear Tool');
  console.log('Uncomment dòng function bạn muốn chạy ở cuối file');
  
  // Uncomment dòng bên dưới để clear toàn bộ cache
  clearAllCache();
}