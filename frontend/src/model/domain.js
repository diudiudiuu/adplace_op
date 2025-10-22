import Base from './base.js';

class Domain extends Base {
    constructor() {
        // 表名
        const tableName = 'tb_domain';

        // 主键
        const primaryKey = 'id';

        // 显示的字段
        const fields = [
            'id',
            'domain',
            'expired_at',
            'ping_flag',
            'hosting',
            'ssl_type',
            "memo",
            "created_at",
            "updated_at"
        ];

        // 字段类型定义
        const fieldTypes = {
            id: {
                type: 'int',
                value: '',
            },
            domain: {
                type: 'string',
                value: '',
            },
            expired_at: {
                type: 'datetime',
                value: '',
            },
            ping_flag: {
                type: 'enum',
                value: ['deploy', 'error', 'normal'],
                default: 'deploy',
            },
            hosting: {
                type: 'enum',
                value: ['cloudflare', 'self'],
                default: 'cloudflare',
            },
            ssl_type: {
                type: 'enum',
                value: ['none', 'letsencrypt'],
                default: 'none',
            },
            memo: {
                type: 'string',
                value: '',
            },
            created_at: {
                type: 'datetime',
                value: '',
            },
            updated_at: {
                type: 'datetime',
                value: '',
            }
        };

        super(tableName, primaryKey, fields, fieldTypes);
        
        // 可选：设置域名数量限制（示例）
        // this.maxRecords = 50;
        // this.maxRecordsMessage = '每个客户最多只能添加50个域名';
    }
}

export default Domain;
