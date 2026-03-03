/**
 * Modern Dashboard Analytics Module
 * Provides interactive charts and visualizations for phishing campaign analytics
 */

// Initialize global namespace
var ModernDashboard = ModernDashboard || {};

// Auto-initialize when DOM is ready
document.addEventListener('DOMContentLoaded', function() {
    if (typeof ModernDashboard !== 'undefined') {
        ModernDashboard.init();
    }
});

// API Configuration
ModernDashboard.api = {
    baseUrl: '/api',
    
    getDashboard: function() {
        return fetch(this.baseUrl + '/analytics/dashboard')
            .then(response => response.json());
    },
    
    getCampaignAnalytics: function(campaignId) {
        return fetch(this.baseUrl + '/analytics/campaign/' + campaignId)
            .then(response => response.json());
    },
    
    getUserScores: function() {
        return fetch(this.baseUrl + '/analytics/users/scores')
            .then(response => response.json());
    },
    
    getDepartmentMetrics: function() {
        return fetch(this.baseUrl + '/analytics/departments')
            .then(response => response.json());
    },
    
    getRiskDistribution: function() {
        return fetch(this.baseUrl + '/analytics/risk/distribution')
            .then(response => response.json());
    },
    
    getTopAtRisk: function(limit) {
        return fetch(this.baseUrl + '/analytics/users/at-risk?limit=' + (limit || 10))
            .then(response => response.json());
    },
    
    getTimeAnalytics: function(campaignId) {
        return fetch(this.baseUrl + '/analytics/campaign/' + campaignId + '/time')
            .then(response => response.json());
    }
};

// Chart Configuration
ModernDashboard.charts = {
    colors: {
        primary: '#6c5ce7',
        secondary: '#a29bfe',
        success: '#00b894',
        warning: '#fdcb6e',
        danger: '#d63031',
        info: '#74b9ff',
        dark: '#2d3436',
        light: '#dfe6e9'
    },
    
    // Risk level colors
    riskColors: {
        low: '#00b894',
        medium: '#fdcb6e',
        high: '#e17055',
        critical: '#d63031'
    },
    
    // Default chart options
    defaultOptions: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: {
                position: 'bottom',
                labels: {
                    padding: 20,
                    usePointStyle: true
                }
            }
        }
    }
};

// Initialize Dashboard
ModernDashboard.init = function() {
    console.log('Initializing Modern Dashboard...');
    this.loadDashboardData();
    this.setupAutoRefresh();
};

// Load all dashboard data
ModernDashboard.loadDashboardData = function() {
    var self = this;
    
    // Show loading state
    $('#analytics-loading').show();
    $('#analytics-content').hide();
    
    Promise.all([
        this.api.getDashboard(),
        this.api.getRiskDistribution(),
        this.api.getTopAtRisk(5)
    ]).then(function(results) {
        var dashboard = results[0];
        var riskDist = results[1];
        var topAtRisk = results[2];
        
        self.renderOverviewCards(dashboard);
        self.renderRiskDistributionChart(riskDist);
        self.renderTopAtRiskTable(topAtRisk);
        self.renderDepartmentChart(dashboard.department_stats);
        
        $('#analytics-loading').hide();
        $('#analytics-content').fadeIn();
    }).catch(function(error) {
        console.error('Error loading dashboard:', error);
        $('#analytics-loading').html(
            '<div class="alert alert-danger">Error loading analytics. Please check your connection.</div>'
        );
    });
};

// Render overview cards
ModernDashboard.renderOverviewCards = function(data) {
    // Total Campaigns
    $('#card-total-campaigns .card-value').text(data.total_campaigns || 0);
    $('#card-active-campaigns .card-value').text(data.active_campaigns || 0);
    
    // Email stats
    $('#card-emails-sent .card-value').text(data.total_emails_sent || 0);
    $('#card-total-clicks .card-value').text(data.total_clicks || 0);
    
    // Rates with color coding
    this.renderRateCard('card-click-rate', data.global_click_rate);
    this.renderRateCard('card-open-rate', data.global_open_rate);
    this.renderRateCard('card-report-rate', data.global_report_rate);
};

ModernDashboard.renderRateCard = function(cardId, rate) {
    var rate = rate || 0;
    var $card = $('#' + cardId);
    var $value = $card.find('.card-value');
    var $badge = $card.find('.rate-badge');
    
    $value.text(rate.toFixed(1) + '%');
    
    // Color coding based on rate
    $badge.removeClass('badge-success badge-warning badge-danger');
    if (cardId.includes('click')) {
        if (rate > 50) {
            $badge.addClass('badge-danger').text('HIGH');
        } else if (rate > 25) {
            $badge.addClass('badge-warning').text('MEDIUM');
        } else {
            $badge.addClass('badge-success').text('GOOD');
        }
    }
};

// Render Risk Distribution Donut Chart
ModernDashboard.renderRiskDistributionChart = function(data) {
    var ctx = document.getElementById('risk-distribution-chart');
    if (!ctx) return;
    
    var labels = ['Low Risk', 'Medium Risk', 'High Risk', 'Critical'];
    var values = [
        data.low || 0,
        data.medium || 0,
        data.high || 0,
        data.critical || 0
    ];
    var colors = [
        this.charts.riskColors.low,
        this.charts.riskColors.medium,
        this.charts.riskColors.high,
        this.charts.riskColors.critical
    ];
    
    new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: labels,
            datasets: [{
                data: values,
                backgroundColor: colors,
                borderWidth: 0
            }]
        },
        options: {
            ...this.charts.defaultOptions,
            cutout: '60%'
        }
    });
};

// Render Top At-Risk Users Table
ModernDashboard.renderTopAtRiskTable = function(users) {
    var tbody = $('#at-risk-users-body');
    tbody.empty();
    
    if (!users || users.length === 0) {
        tbody.html('<tr><td colspan="5" class="text-center">No data available</td></tr>');
        return;
    }
    
    users.forEach(function(user, index) {
        var riskClass = 'risk-' + user.risk_level;
        var row = `
            <tr>
                <td>${index + 1}</td>
                <td>${user.email}</td>
                <td><span class="score-badge ${riskClass}">${user.score}</span></td>
                <td>
                    <div class="progress" style="height: 8px;">
                        <div class="progress-bar ${riskClass}" style="width: ${user.score}%"></div>
                    </div>
                </td>
                <td><span class="label label-${user.risk_level === 'critical' ? 'danger' : user.risk_level === 'high' ? 'warning' : 'success'}">${user.risk_level}</span></td>
            </tr>
        `;
        tbody.append(row);
    });
};

// Render Department Chart
ModernDashboard.renderDepartmentChart = function(departments) {
    var ctx = document.getElementById('department-chart');
    if (!ctx || !departments || departments.length === 0) return;
    
    var labels = departments.map(d => d.department);
    var clickRates = departments.map(d => d.click_rate);
    var reportRates = departments.map(d => d.report_rate);
    
    new Chart(ctx, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [
                {
                    label: 'Click Rate %',
                    data: clickRates,
                    backgroundColor: this.charts.colors.danger,
                    borderRadius: 4
                },
                {
                    label: 'Report Rate %',
                    data: reportRates,
                    backgroundColor: this.charts.colors.success,
                    borderRadius: 4
                }
            ]
        },
        options: {
            ...this.charts.defaultOptions,
            scales: {
                y: {
                    beginAtZero: true,
                    max: 100
                }
            }
        }
    });
};

// Campaign-specific analytics
ModernDashboard.loadCampaignAnalytics = function(campaignId) {
    var self = this;
    
    Promise.all([
        this.api.getCampaignAnalytics(campaignId),
        this.api.getTimeAnalytics(campaignId)
    ]).then(function(results) {
        var analytics = results[0];
        var timeAnalytics = results[1];
        
        // Update campaign stats
        $('#campaign-sent').text(analytics.emails_sent);
        $('#campaign-opened').text(analytics.emails_opened);
        $('#campaign-clicked').text(analytics.links_clicked);
        $('#campaign-reported').text(analytics.emails_reported);
        
        // Update rates
        $('#campaign-click-rate').text((analytics.click_rate || 0).toFixed(1) + '%');
        $('#campaign-open-rate').text((analytics.open_rate || 0).toFixed(1) + '%');
        $('#campaign-report-rate').text((analytics.report_rate || 0).toFixed(1) + '%');
        
        // Render time distribution chart
        self.renderTimeDistributionChart(timeAnalytics);
        
    }).catch(function(error) {
        console.error('Error loading campaign analytics:', error);
    });
};

ModernDashboard.renderTimeDistributionChart = function(data) {
    var ctx = document.getElementById('time-distribution-chart');
    if (!ctx) return;
    
    new Chart(ctx, {
        type: 'bar',
        data: {
            labels: ['< 1 min', '1-5 min', '5-30 min', '30-60 min', '> 1 hour'],
            datasets: [{
                label: 'Clicks',
                data: [
                    data.clicks_within_1min || 0,
                    data.clicks_1_to_5min || 0,
                    data.clicks_5_to_30min || 0,
                    data.clicks_30_to_60min || 0,
                    data.clicks_after_1hour || 0
                ],
                backgroundColor: this.charts.colors.primary,
                borderRadius: 4
            }]
        },
        options: {
            ...this.charts.defaultOptions,
            plugins: {
                title: {
                    display: true,
                    text: 'Time to Click Distribution'
                }
            }
        }
    });
};

// Auto-refresh every 30 seconds
ModernDashboard.setupAutoRefresh = function() {
    setInterval(() => {
        this.loadDashboardData();
    }, 30000);
};

// Export for use
if (typeof module !== 'undefined' && module.exports) {
    module.exports = ModernDashboard;
}
