using Microsoft.EntityFrameworkCore;
using SimpleTodoAPI.Models;

namespace SimpleTodoAPI.Data
{
    public class TodoContext : DbContext
    {
        public TodoContext(DbContextOptions<TodoContext> options) : base(options)
        {
        }
        
        public DbSet<Todo> Todos { get; set; }
    }
}