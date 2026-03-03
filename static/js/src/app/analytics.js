var Dashboard = {
    riskChart: null,
    campaignChart: null,
    
    init: function() {
        this.load();
    },
    
    load: function() {
        console.log("Loading analytics...");
        console.log("User API key:", user ? user.api_key : "NOT SET");
        $("#loading").show();
        $("#dashboard-content").hide();
        
        api.analytics.summary().then(function(response) {
            console.log("Summary response:", response);
            Dashboard.renderSummary(response);
            $("#loading").hide();
            $("#dashboard-content").show();
        }).catch(function(err) {
            console.error("Error loading summary:", err);
            $("#loading").hide();
            $("#dashboard-content").show();
            $("#dashboard-content").html("<div class='alert alert-warning'>Please login to view analytics</div>");
        });
        
        api.analytics.departments().then(function(response) {
            console.log("Departments response:", response);
            Dashboard.renderDepartments(response);
        }).catch(function(err) {
            console.error("Error loading departments:", err);
        });
        
        api.analytics.userScores().then(function(response) {
            console.log("User scores response:", response);
            Dashboard.renderUserScores(response);
        }).catch(function(err) {
            console.error("Error loading user scores:", err);
        });
        
        api.analytics.campaigns().then(function(response) {
            console.log("Campaigns response:", response);
            Dashboard.renderCampaigns(response);
        }).catch(function(err) {
            console.error("Error loading campaigns:", err);
        });
    },
    
    renderSummary: function(data) {
        $("#stat-campaigns").text(data.total_campaigns || 0);
        $("#stat-emails").text(data.total_emails_sent || 0);
        $("#stat-click-rate").text((data.overall_click_rate || 0).toFixed(1) + "%");
        $("#stat-report-rate").text((data.overall_report_rate || 0).toFixed(1) + "%");
        
        Dashboard.renderRiskChart(data.risk_distribution);
    },
    
    renderRiskChart: function(distribution) {
        var ctx = document.getElementById('riskChart');
        if (!ctx) return;
        
        if (Dashboard.riskChart) {
            Dashboard.riskChart.destroy();
        }
        
        var data = {
            labels: ['Low', 'Medium', 'High', 'Critical'],
            datasets: [{
                data: [
                    distribution.low || 0,
                    distribution.medium || 0,
                    distribution.high || 0,
                    distribution.critical || 0
                ],
                backgroundColor: [
                    '#28a745',
                    '#ffc107',
                    '#fd7e14',
                    '#dc3545'
                ]
            }]
        };
        
        Dashboard.riskChart = new Chart(ctx, {
            type: 'doughnut',
            data: data,
            options: {
                responsive: true,
                plugins: {
                    legend: {
                        position: 'bottom'
                    }
                }
            }
        });
    },
    
    renderCampaigns: function(campaigns) {
        var ctx = document.getElementById('campaignChart');
        if (!ctx) return;
        
        if (Dashboard.campaignChart) {
            Dashboard.campaignChart.destroy();
        }
        
        var campaignData = campaigns.slice(0, 5).reverse();
        
        Dashboard.campaignChart = new Chart(ctx, {
            type: 'bar',
            data: {
                labels: campaignData.map(function(c) { return c.campaign_name || 'Campaign ' + c.campaign_id; }),
                datasets: [
                    {
                        label: 'Sent',
                        data: campaignData.map(function(c) { return c.emails_sent || 0; }),
                        backgroundColor: '#28a745'
                    },
                    {
                        label: 'Opened',
                        data: campaignData.map(function(c) { return c.emails_opened || 0; }),
                        backgroundColor: '#17a2b8'
                    },
                    {
                        label: 'Clicked',
                        data: campaignData.map(function(c) { return c.links_clicked || 0; }),
                        backgroundColor: '#ffc107'
                    },
                    {
                        label: 'Reported',
                        data: campaignData.map(function(c) { return c.emails_reported || 0; }),
                        backgroundColor: '#dc3545'
                    }
                ]
            },
            options: {
                responsive: true,
                scales: {
                    x: {
                        stacked: false
                    },
                    y: {
                        beginAtZero: true
                    }
                }
            }
        });
    },
    
    renderDepartments: function(departments) {
        var tbody = $("#departmentBody");
        tbody.empty();
        
        if (!departments || departments.length === 0) {
            tbody.html('<tr><td colspan="7" class="text-center">No department data available</td></tr>');
            return;
        }
        
        departments.forEach(function(dept) {
            var riskClass = 'success';
            if (dept.high_risk_users > dept.medium_risk_users) {
                riskClass = 'danger';
            } else if (dept.medium_risk_users > 0) {
                riskClass = 'warning';
            }
            
            var row = '<tr>' +
                '<td>' + (dept.department || 'N/A') + '</td>' +
                '<td>' + (dept.domain || 'N/A') + '</td>' +
                '<td>' + dept.total_users + '</td>' +
                '<td>' + dept.click_rate.toFixed(1) + '%</td>' +
                '<td>' + dept.open_rate.toFixed(1) + '%</td>' +
                '<td>' + dept.report_rate.toFixed(1) + '%</td>' +
                '<td><span class="label label-' + riskClass + '">' + 
                (dept.high_risk_users > 0 ? 'High' : (dept.medium_risk_users > 0 ? 'Medium' : 'Low')) + 
                '</span></td>' +
                '</tr>';
            tbody.append(row);
        });
    },
    
    renderUserScores: function(scores) {
        var tbody = $("#userScoresBody");
        tbody.empty();
        
        if (!scores || scores.length === 0) {
            tbody.html('<tr><td colspan="7" class="text-center">No user scores available. Run some campaigns to see data.</td></tr>');
            return;
        }
        
        scores.slice(0, 20).forEach(function(score) {
            var riskClass = 'success';
            var riskLabel = 'Low';
            
            switch(score.risk_level) {
                case 'critical':
                    riskClass = 'danger';
                    riskLabel = 'Critical';
                    break;
                case 'high':
                    riskClass = 'danger';
                    riskLabel = 'High';
                    break;
                case 'medium':
                    riskClass = 'warning';
                    riskLabel = 'Medium';
                    break;
            }
            
            var row = '<tr>' +
                '<td>' + score.email + '</td>' +
                '<td>' + (score.department || 'N/A') + '</td>' +
                '<td>' + score.times_clicked + '</td>' +
                '<td>' + score.times_opened + '</td>' +
                '<td>' + score.times_reported + '</td>' +
                '<td><strong>' + score.score + '</strong></td>' +
                '<td><span class="label label-' + riskClass + '">' + riskLabel + '</span></td>' +
                '</tr>';
            tbody.append(row);
        });
    }
};

$(document).ready(function() {
    if ($("#riskChart").length > 0) {
        Dashboard.init();
    }
});
