const std = @import("std");
const config = @import("config.zig");
const domain = @import("domain.zig");
pub const TmuxSession = domain.TmuxSession;
pub const TmuxEnvironment = domain.TmuxEnvironment;
pub const TmuxPane = domain.TmuxPane;

pub fn getSessions(allocator: std.mem.Allocator) ![]const TmuxSession {
    const cfg = try config.fetchConfig(allocator);
    return cfg.entries;
}

pub fn appendSession(allocator: std.mem.Allocator, session: TmuxSession) !void {
    const sessions = try getSessions(allocator);
    var new_sessions = std.ArrayList(TmuxSession).init(allocator);
    for (sessions) |entry| {
        if (std.mem.eql(u8, entry.name, session.name)) {
            std.debug.print("Already exists: {s}\n", .{entry.name});
            return;
        }
        _ = try new_sessions.append(entry);
    }
    _ = try new_sessions.append(session);
    try config.saveConfig(allocator, config.Config{ .entries = try new_sessions.toOwnedSlice() });
}

pub fn prepareSession(key: []const u8, allocator: std.mem.Allocator) !void {
    _ = key;
    _ = allocator;
    return;
}

pub fn deleteSession(key: []const u8, allocator: std.mem.Allocator) !void {
    const sessions = try getSessions(allocator);
    var new_sessions = std.ArrayList(TmuxSession).init(allocator);
    for (sessions) |entry| {
        if (!std.mem.eql(u8, try sessionToKey(allocator, entry), key)) {
            _ = try new_sessions.append(entry);
        } else {
            std.debug.print("Removing: {s}\n", .{entry.name});
        }
    }
    _ = try config.saveConfig(allocator, config.Config{ .entries = try new_sessions.toOwnedSlice() });
}

pub fn sessionToKey(allocator: std.mem.Allocator, entry: TmuxSession) ![]const u8 {
    if (std.mem.eql(u8, entry.name, entry.path)) {
        return entry.name;
    } else {
        const result = try std.fmt.allocPrint(allocator, "{s} ({s})", .{ entry.name, entry.path });
        return result;
    }
}
