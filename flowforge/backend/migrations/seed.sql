-- Seed Data for FlowForge
-- Run this after migrations to populate initial data

-- Note: Password for all users is 'password123' (hashed with bcrypt)
-- Hash: $2a$10$rZ5q2WX8oT5QzYxKj9ZvKe3YnF8YnKxYnZxYnZxYnZxYnZxYnZx

-- Create a demo organization
INSERT INTO organizations (name) VALUES ('Demo Organization');

-- Create users (password: password123 for all)
INSERT INTO users (organization_id, email, password_hash, name, role) VALUES 
(1, 'admin@demo.com', '$2a$10$rZ5q2WX8oT5QzYxKj9ZvKe3YnF8YnKxYnZxYnZxYnZxYnZxYnZx', 'Admin User', 'admin'),
(1, 'manager@demo.com', '$2a$10$rZ5q2WX8oT5QzYxKj9ZvKe3YnF8YnKxYnZxYnZxYnZxYnZxYnZx', 'Manager User', 'manager'),
(1, 'operator@demo.com', '$2a$10$rZ5q2WX8oT5QzYxKj9ZvKe3YnF8YnKxYnZxYnZxYnZxYnZxYnZx', 'Operator User', 'operator'),
(1, 'viewer@demo.com', '$2a$10$rZ5q2WX8oT5QzYxKj9ZvKe3YnF8YnKxYnZxYnZxYnZxYnZxYnZx', 'Viewer User', 'viewer');

-- Create a workflow template
INSERT INTO workflow_templates (organization_id, name, description, version, is_active, created_by) VALUES 
(1, 'Expense Approval', 'Standard expense approval workflow', 1, 1, 1);

-- Create template stages
INSERT INTO template_stages (template_id, name, order_index, required_role, approval_type, timeout_minutes, escalation_role) VALUES 
(1, 'Submit Expense', 0, 'operator', 'manual', 0, NULL),
(1, 'Manager Review', 1, 'manager', 'manual', 1440, 'admin'),
(1, 'Finance Approval', 2, 'admin', 'manual', 2880, NULL),
(1, 'Completed', 3, 'viewer', 'auto', 0, NULL);

-- Create another template
INSERT INTO workflow_templates (organization_id, name, description, version, is_active, created_by) VALUES 
(1, 'Leave Request', 'Employee leave approval workflow', 1, 1, 1);

INSERT INTO template_stages (template_id, name, order_index, required_role, approval_type, timeout_minutes, escalation_role) VALUES 
(2, 'Submit Request', 0, 'operator', 'manual', 0, NULL),
(2, 'Manager Approval', 1, 'manager', 'manual', 720, 'admin'),
(2, 'HR Review', 2, 'admin', 'manual', 1440, NULL);
