-- ============================================================
-- AI 学习平台 - 用户种子数据
-- 创建管理员/讲师/学员账户（各 3 个）
-- ============================================================

-- 管理员账户 (3 个)
INSERT INTO users (id, username, email, password_hash, avatar_url, role, is_active, created_at, last_login_at) VALUES
('00000000-0000-0000-0000-000000000001', 'admin_zhang', 'zhang.admin@aiplatform.com', '$2a$10$X7K9ZjYqN5xL8vM3wR2tUeH6fG4pQ1sW9nK0mL7cV8bD5aE3fH2iJ', 'https://api.dicebear.com/7.x/avataaars/svg?seed=admin_zhang', 'admin', true, NOW() - INTERVAL '30 days', NOW() - INTERVAL '1 day'),
('00000000-0000-0000-0000-000000000002', 'admin_wang', 'wang.admin@aiplatform.com', '$2a$10$X7K9ZjYqN5xL8vM3wR2tUeH6fG4pQ1sW9nK0mL7cV8bD5aE3fH2iJ', 'https://api.dicebear.com/7.x/avataaars/svg?seed=admin_wang', 'admin', true, NOW() - INTERVAL '25 days', NOW() - INTERVAL '2 days'),
('00000000-0000-0000-0000-000000000003', 'admin_li', 'li.admin@aiplatform.com', '$2a$10$X7K9ZjYqN5xL8vM3wR2tUeH6fG4pQ1sW9nK0mL7cV8bD5aE3fH2iJ', 'https://api.dicebear.com/7.x/avataaars/svg?seed=admin_li', 'admin', true, NOW() - INTERVAL '20 days', NOW() - INTERVAL '3 days');

-- 讲师账户 (3 个)
INSERT INTO users (id, username, email, password_hash, avatar_url, role, is_active, created_at, last_login_at) VALUES
('10000000-0000-0000-0000-000000000001', 'prof_chen', 'chen.professor@aiplatform.com', '$2a$10$X7K9ZjYqN5xL8vM3wR2tUeH6fG4pQ1sW9nK0mL7cV8bD5aE3fH2iJ', 'https://api.dicebear.com/7.x/avataaars/svg?seed=prof_chen', 'instructor', true, NOW() - INTERVAL '60 days', NOW() - INTERVAL '1 hour'),
('10000000-0000-0000-0000-000000000002', 'prof_liu', 'liu.professor@aiplatform.com', '$2a$10$X7K9ZjYqN5xL8vM3wR2tUeH6fG4pQ1sW9nK0mL7cV8bD5aE3fH2iJ', 'https://api.dicebear.com/7.x/avataaars/svg?seed=prof_liu', 'instructor', true, NOW() - INTERVAL '55 days', NOW() - INTERVAL '2 hours'),
('10000000-0000-0000-0000-000000000003', 'prof_huang', 'huang.professor@aiplatform.com', '$2a$10$X7K9ZjYqN5xL8vM3wR2tUeH6fG4pQ1sW9nK0mL7cV8bD5aE3fH2iJ', 'https://api.dicebear.com/7.x/avataaars/svg?seed=prof_huang', 'instructor', true, NOW() - INTERVAL '50 days', NOW() - INTERVAL '3 hours');

-- 学员账户 (3 个)
INSERT INTO users (id, username, email, password_hash, avatar_url, role, is_active, created_at, last_login_at) VALUES
('20000000-0000-0000-0000-000000000001', 'student_zhao', 'zhao.student@aiplatform.com', '$2a$10$X7K9ZjYqN5xL8vM3wR2tUeH6fG4pQ1sW9nK0mL7cV8bD5aE3fH2iJ', 'https://api.dicebear.com/7.x/avataaars/svg?seed=student_zhao', 'student', true, NOW() - INTERVAL '15 days', NOW() - INTERVAL '30 minutes'),
('20000000-0000-0000-0000-000000000002', 'student_wu', 'wu.student@aiplatform.com', '$2a$10$X7K9ZjYqN5xL8vM3wR2tUeH6fG4pQ1sW9nK0mL7cV8bD5aE3fH2iJ', 'https://api.dicebear.com/7.x/avataaars/svg?seed=student_wu', 'student', true, NOW() - INTERVAL '14 days', NOW() - INTERVAL '1 hour'),
('20000000-0000-0000-0000-000000000003', 'student_zhou', 'zhou.student@aiplatform.com', '$2a$10$X7K9ZjYqN5xL8vM3wR2tUeH6fG4pQ1sW9nK0mL7cV8bD5aE3fH2iJ', 'https://api.dicebear.com/7.x/avataaars/svg?seed=student_zhou', 'student', true, NOW() - INTERVAL '13 days', NOW() - INTERVAL '2 hours');

-- 密码说明：所有账户默认密码为 "Password123!"
-- 生产环境请修改密码哈希值
