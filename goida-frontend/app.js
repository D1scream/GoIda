const { createApp, reactive } = Vue;
const api = {
    baseURL: 'http://localhost:8080/api',
    async request(endpoint, options = {}) {
        const config = { headers: { 'Content-Type': 'application/json', ...options.headers }, ...options };
        const token = localStorage.getItem('authToken');
        if (token) config.headers['Authorization'] = `Bearer ${token}`;
        try {
            const response = await axios({ url: `${this.baseURL}${endpoint}`, ...config });
            return response.data;
        } catch (error) {
            throw new Error(error.response?.data?.message || error.message);
        }
    },
    async get(endpoint) { return this.request(endpoint); },
    async post(endpoint, data) { return this.request(endpoint, { method: 'POST', data }); },
    async put(endpoint, data) { return this.request(endpoint, { method: 'PUT', data }); },
    async delete(endpoint) { return this.request(endpoint, { method: 'DELETE' }); }
};

const ArticleCard = {
    props: ['article', 'currentUser'],
    emits: ['delete', 'edit'],
    template: `<div class="article-card">
        <h4>{{ article.title }}</h4>
        <p><strong>Автор:</strong> {{ article.author_name }}</p>
        <p><strong>Создана:</strong> {{ formatDate(article.created_at) }}</p>
        <p><strong>Содержимое:</strong> {{ truncateText(article.content, 200) }}</p>
        <div class="article-actions">
            <button v-if="canEdit" class="btn-small btn-warning" @click="$emit('edit', article)">Редактировать</button>
            <button v-if="canEdit" class="btn-small btn-danger" @click="confirmDelete">Удалить</button>
        </div>
    </div>`,
    computed: {
        canEdit() {
            return this.currentUser && (this.currentUser.role.name === 'admin' || this.currentUser.id === this.article.author_id);
        }
    },
    methods: {
        formatDate(dateString) { return new Date(dateString).toLocaleString('ru-RU'); },
        truncateText(text, maxLength) { return text.length > maxLength ? text.substring(0, maxLength) + '...' : text; },
        confirmDelete() { if (confirm('Вы уверены, что хотите удалить эту статью?')) this.$emit('delete', this.article.id); }
    }
};

const UserCard = {
    props: ['user'],
    template: `<div class="user-card">
        <h4>{{ user.name }}</h4>
        <p><strong>Email:</strong> {{ user.email }}</p>
        <p><strong>Роль:</strong> {{ user.role.name }}</p>
        <p><strong>Зарегистрирован:</strong> {{ formatDate(user.created_at) }}</p>
    </div>`,
    methods: {
        formatDate(dateString) { return new Date(dateString).toLocaleString('ru-RU'); }
    }
};

const appData = reactive({
    loginForm: { login: 'admin', password: 'password' },
    registerForm: { name: '', email: '', login: '', password: '', confirmPassword: '' },
    isAuthenticated: false, currentUser: null, authToken: null, showRegisterForm: false,
    isLoading: false, authStatus: null, articleForm: { title: '', content: '' },
    editForm: { id: null, title: '', content: '' }, showEditModal: false,
    articles: [], users: [], logs: []
});

const app = createApp({
    components: { ArticleCard, UserCard },
    data() { return appData; },
    computed: {
        isAdmin() { return this.currentUser && this.currentUser.role && this.currentUser.role.name === 'admin'; }
    },
    methods: {
        async login() {
            if (!this.loginForm.login || !this.loginForm.password) {
                this.showStatus('Пожалуйста, введите логин и пароль', 'error');
                return;
            }
            this.isLoading = true;
            try {
                const response = await api.post('/auth/login', { login: this.loginForm.login, password: this.loginForm.password });
                this.authToken = response.token;
                this.currentUser = response.user;
                this.isAuthenticated = true;
                localStorage.setItem('authToken', this.authToken);
                localStorage.setItem('currentUser', JSON.stringify(this.currentUser));
                this.showStatus(`Добро пожаловать, ${this.currentUser.name}!`, 'success');
                this.addLog('Успешная авторизация', 'success');
            } catch (error) {
                this.showStatus(`Ошибка авторизации: ${error.message}`, 'error');
                this.addLog(`Ошибка авторизации: ${error.message}`, 'error');
            } finally {
                this.isLoading = false;
            }
        },

        logout() {
            this.authToken = null; this.currentUser = null; this.isAuthenticated = false;
            localStorage.removeItem('authToken'); localStorage.removeItem('currentUser');
            this.showStatus('Вы вышли из системы', 'info');
            this.articles = []; this.users = []; this.addLog('Пользователь вышел из системы', 'info');
        },

        async register() {
            if (!this.registerForm.name || !this.registerForm.email || !this.registerForm.login || !this.registerForm.password) {
                this.showStatus('Пожалуйста, заполните все поля', 'error'); return;
            }
            if (this.registerForm.password !== this.registerForm.confirmPassword) {
                this.showStatus('Пароли не совпадают', 'error'); return;
            }
            if (this.registerForm.password.length < 6) {
                this.showStatus('Пароль должен содержать минимум 6 символов', 'error'); return;
            }
            this.isLoading = true;
            try {
                await api.post('/users', { name: this.registerForm.name, email: this.registerForm.email, login: this.registerForm.login, password: this.registerForm.password });
                this.showStatus('Регистрация успешна! Теперь вы можете войти в систему', 'success');
                this.addLog('Пользователь зарегистрирован', 'success');
                this.registerForm = { name: '', email: '', login: '', password: '', confirmPassword: '' };
                this.showRegisterForm = false;
            } catch (error) {
                this.showStatus(`Ошибка регистрации: ${error.message}`, 'error');
                this.addLog(`Ошибка регистрации: ${error.message}`, 'error');
            } finally {
                this.isLoading = false;
            }
        },


        async createArticle() {
            if (!this.isAuthenticated) { this.showStatus('Необходимо авторизоваться', 'error'); return; }
            if (!this.articleForm.title || !this.articleForm.content) { this.showStatus('Пожалуйста, заполните все поля', 'error'); return; }
            this.isLoading = true;
            try {
                await api.post('/articles', { title: this.articleForm.title, content: this.articleForm.content });
                this.showStatus('Статья успешно создана!', 'success');
                this.articleForm = { title: '', content: '' };
                this.loadArticles(); this.addLog('Статья создана', 'success');
            } catch (error) {
                this.showStatus(`Ошибка создания статьи: ${error.message}`, 'error');
                this.addLog(`Ошибка создания статьи: ${error.message}`, 'error');
            } finally { this.isLoading = false; }
        },

        async loadArticles() {
            if (!this.isAuthenticated) { this.showStatus('Необходимо авторизоваться', 'error'); return; }
            this.isLoading = true;
            try {
                this.articles = await api.get('/articles?limit=50');
                this.addLog(`Загружено ${this.articles.length} статей`, 'success');
            } catch (error) {
                this.showStatus(`Ошибка загрузки статей: ${error.message}`, 'error');
                this.addLog(`Ошибка загрузки статей: ${error.message}`, 'error');
            } finally { this.isLoading = false; }
        },

        async deleteArticle(articleId) {
            if (!this.isAuthenticated) { this.showStatus('Необходимо авторизоваться', 'error'); return; }
            this.isLoading = true;
            try {
                await api.delete(`/articles/${articleId}`);
                this.showStatus('Статья удалена', 'success');
                this.loadArticles(); this.addLog(`Статья ${articleId} удалена`, 'success');
            } catch (error) {
                this.showStatus(`Ошибка удаления статьи: ${error.message}`, 'error');
                this.addLog(`Ошибка удаления статьи: ${error.message}`, 'error');
            } finally { this.isLoading = false; }
        },
        editArticle(article) { alert(`Редактирование статьи: ${article.title}`); },

        async loadUsers() {
            if (!this.isAuthenticated) { this.showStatus('Необходимо авторизоваться', 'error'); return; }
            if (!this.isAdmin) { this.showStatus('Доступ только для администраторов', 'error'); return; }
            this.isLoading = true;
            try {
                this.users = await api.get('/admin/users?limit=50');
                this.addLog(`Загружено ${this.users.length} пользователей`, 'success');
            } catch (error) {
                this.showStatus(`Ошибка загрузки пользователей: ${error.message}`, 'error');
                this.addLog(`Ошибка загрузки пользователей: ${error.message}`, 'error');
            } finally { this.isLoading = false; }
        },

        editArticle(article) {
            this.editForm = { id: article.id, title: article.title, content: article.content };
            this.showEditModal = true;
        },

        closeEditModal() {
            this.showEditModal = false;
            this.editForm = { id: null, title: '', content: '' };
        },

        async updateArticle() {
            if (!this.editForm.title || !this.editForm.content) {
                this.showStatus('Пожалуйста, заполните все поля', 'error');
                return;
            }
            this.isLoading = true;
            try {
                await api.put(`/articles/${this.editForm.id}`, { title: this.editForm.title, content: this.editForm.content });
                this.showStatus('Статья успешно обновлена!', 'success');
                this.closeEditModal();
                this.loadArticles();
                this.addLog('Статья обновлена', 'success');
            } catch (error) {
                this.showStatus(`Ошибка обновления статьи: ${error.message}`, 'error');
                this.addLog(`Ошибка обновления статьи: ${error.message}`, 'error');
            } finally {
                this.isLoading = false;
            }
        },

        showStatus(message, type) { this.authStatus = { message, type }; setTimeout(() => { this.authStatus = null; }, 5000); },
        addLog(message, type) {
            const timestamp = new Date().toLocaleTimeString('ru-RU');
            this.logs.push({ message: `[${timestamp}] ${message}`, type: `log-${type}` });
            if (this.logs.length > 100) this.logs.shift();
        },
        clearLogs() { this.logs = []; }
    },

    mounted() {
        const savedUser = localStorage.getItem('currentUser');
        const savedToken = localStorage.getItem('authToken');
        if (savedUser && savedToken) {
            this.currentUser = JSON.parse(savedUser);
            this.authToken = savedToken;
            this.isAuthenticated = true;
        }
        this.addLog('Vue.js фронтенд загружен и готов к работе', 'info');
    }
});

app.mount('#app');
